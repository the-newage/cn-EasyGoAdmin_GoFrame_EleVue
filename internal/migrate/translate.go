// +======================================================================
// | EasyGoAdmin 企业数字化开发平台 [ 赋能开发者，助力企业发展 ]
// +----------------------------------------------------------------------
// | Copyright (c) 2019-2026 半城风雨团队 版权所有
// +----------------------------------------------------------------------
// | 开源协议: Apache-2.0
// | 重要声明: 未经书面许可，不得移除或修改本框架相关版权标识。
// +----------------------------------------------------------------------
// | 官方网站: https://www.easygoadmin.vip
// +----------------------------------------------------------------------
// | 研发团队: 小蚂蚁云团队 荣誉出品
// +======================================================================
// | 版权和免责声明:
// | 一、知识产权声明
// | 本团队对该软件框架产品（以下简称“本产品”）拥有完整知识产权，包括但不限于商标权、专利权、
// | 著作权、商业秘密等。前述权利均受中华人民共和国及相关国际法律法规保护。未经本团队书面授权，
// | 任何个人或单位不得基于本产品本身申请任何形式的知识产权。
// | 二、使用限制与禁止条款
// | 您理解并同意，禁止将本产品用于以下任何行为：
// | 1. 任何违反中华人民共和国法律法规的用途；
// | 2. 任何侵害他人合法权益（包括但不限于人身权、财产权、知识产权等）的行为；
// | 3. 任何危害国家安全、公共利益及社会道德风尚的行为。
// | 三、责任豁免
// | 因使用本产品而产生的任何损失或法律纠纷（包括但不限于直接损失、间接损失、附带损失、衍生损失
// | 以及因违约、侵权等引发的责任），本团队均不承担任何法律责任。若因违规使用本产品导致任何第三
// | 方提起索赔或行政机关进行调查，您应独立承担全部责任。
// | 四、违约与法律责任
// | 如发现任何违反本声明的行为，本团队有权：
// | 1. 无条件配合公安机关及其他有权机关的调查取证工作；
// | 2. 依法追究相关方的民事责任，并要求赔偿全部损失；
// | 3. 保留进一步采取刑事控告等法律措施的权利。
// | 五、许可范围说明
// | 1、被授权方应视本产品源代码为核心机密信息，负有最高标准的保密与保管义务。禁止以任何形式（包
// | 括但不限于直接或间接的分享、复制、传播、分发、转售、出租、开源、反向工程、反编译、反汇编或通
// | 过网络服务提供给任何第三方）泄露、披露或处置源代码的全部或部分内容。任何未经授权的复制、使用、
// | 传播行为均构成对本团队知识产权的严重侵权，本团队有权依法追究其全部法律责任，并要求赔偿因此造
// | 成的一切经济损失及商誉损害。
// | 2、本产品仅限于公司或个人在法律允许范围内，用于内部且合规的软件产品研发。具体权利义务以《软件
// | 许可协议》附件为准。
// +======================================================================

package migrate

import (
	"regexp"
	"strconv"
	"strings"
)

// ============================== 类型归一 ==============================

var (
	reParen    = regexp.MustCompile(`\(([^)]*)\)`)               // (11) / (10,2) / ('a','b')
	reUnsigned = regexp.MustCompile(`(?i)\bunsigned\b`)          // mysql unsigned 修饰
	rePgCast   = regexp.MustCompile(`(?i)^'(.*)'::[\w\s\[\]]+$`) // pgsql 'x'::character varying
	rePgFunc   = regexp.MustCompile(`(?i)^(current_timestamp|now|localtime)\(\d*\)$`)
)

// normalizeType 源方言原始列类型转换为中间类型
//
// 同时解析长度、精度与枚举值
func normalizeType(dialect, rawType string) (typ string, length, scale int, enumVals []string) {
	t := strings.ToLower(strings.TrimSpace(rawType))
	unsigned := reUnsigned.MatchString(t)
	t = reUnsigned.ReplaceAllString(t, "")
	t = strings.TrimSpace(t)

	// 提取括号参数（长度/精度/枚举值）
	pargs := ""
	if m := reParen.FindStringSubmatch(t); m != nil {
		pargs = strings.TrimSpace(m[1])
		t = strings.TrimSpace(reParen.ReplaceAllString(t, ""))
	}
	// 剥离 pgsql 时区/精度修饰
	t = strings.TrimSuffix(t, " without time zone")
	t = strings.TrimSuffix(t, " with time zone")

	// 枚举/集合：括号内为引号列表
	if strings.HasPrefix(t, "enum") || strings.HasPrefix(t, "set") || strings.HasPrefix(t, "enum(") {
		for _, v := range strings.Split(pargs, ",") {
			if v = strings.Trim(strings.TrimSpace(v), "'\""); v != "" {
				enumVals = append(enumVals, v)
			}
		}
		if strings.HasPrefix(t, "set") {
			return TypeSet, 0, 0, enumVals
		}
		return TypeEnum, 0, 0, enumVals
	}

	// decimal/numeric：pargs = "10,2"
	if parts := strings.Split(pargs, ","); len(parts) == 2 {
		if p, e1 := strconv.Atoi(strings.TrimSpace(parts[0])); e1 == nil {
			if s, e2 := strconv.Atoi(strings.TrimSpace(parts[1])); e2 == nil {
				return TypeDecimal, p, s, nil
			}
		}
	}

	// 单参数括号：varchar(150) / int(11)
	if pargs != "" {
		if n, err := strconv.Atoi(strings.TrimSpace(strings.Split(pargs, ",")[0])); err == nil {
			length = n
		}
	}

	switch t {
	case "bool", "boolean", "bit":
		return TypeBool, 0, 0, nil
	case "tinyint":
		// mysql tinyint(1) 惯例为布尔
		if length == 1 {
			return TypeBool, 0, 0, nil
		}
		if unsigned {
			return TypeSmallint, 0, 0, nil
		}
		return TypeTinyint, 0, 0, nil
	case "smallint", "year":
		return TypeSmallint, 0, 0, nil
	case "mediumint", "int", "integer":
		return TypeInt, 0, 0, nil
	case "bigint":
		return TypeBigint, 0, 0, nil
	case "float", "real", "binary_float":
		return TypeFloat, 0, 0, nil
	case "double", "double precision", "float8", "binary_double":
		return TypeDouble, 0, 0, nil
	case "decimal", "numeric", "number", "dec", "fixed":
		if length == 0 {
			return TypeDecimal, 0, 0, nil
		}
		return TypeDecimal, length, scale, nil
	case "varchar", "character varying", "nvarchar", "varchar2":
		return TypeVarchar, length, 0, nil
	case "char", "character", "nchar", "char2":
		return TypeChar, length, 0, nil
	case "tinytext":
		return TypeTinytext, 0, 0, nil
	case "text", "mediumtext", "longtext", "clob", "ntext", "json", "jsonb", "xml":
		if t == "json" || t == "jsonb" {
			return TypeJson, 0, 0, nil
		}
		if t == "mediumtext" {
			return TypeMediumText, 0, 0, nil
		}
		if t == "longtext" {
			return TypeLongText, 0, 0, nil
		}
		return TypeText, 0, 0, nil
	case "date":
		return TypeDate, 0, 0, nil
	case "time":
		return TypeTime, 0, 0, nil
	case "datetime", "timestamp", "smalldatetime", "datetime2":
		return TypeDateTime, 0, 0, nil
	case "tinyblob", "blob", "mediumblob", "longblob", "binary", "varbinary",
		"bytea", "image", "raw", "long raw":
		return TypeBlob, 0, 0, nil
	case "uuid":
		return TypeChar, 36, 0, nil
	}
	return TypeOther, 0, 0, nil
}

// cleanDefault 归一默认值表达式
//
// PostgreSQL 剥去类型转换与函数括号，其余方言原样保留
func cleanDefault(dialect string, raw any) string {
	if raw == nil {
		return ""
	}
	s := strings.TrimSpace(toString(raw))
	if s == "" || strings.EqualFold(s, "null") {
		return ""
	}
	if dialect == "pgsql" {
		if m := rePgCast.FindStringSubmatch(s); m != nil {
			s = m[1]
		} else if rePgFunc.MatchString(s) {
			s = "CURRENT_TIMESTAMP"
		}
	}
	if dialect == "mysql" && strings.EqualFold(s, "current_timestamp()") {
		s = "CURRENT_TIMESTAMP"
	}
	return s
}

// toString 将任意默认值转为字符串
// 兼容不同驱动的返回类型差异
func toString(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case []byte:
		return string(t)
	default:
		return strings.TrimSpace(strings.TrimSuffix(strval(t), ".0"))
	}
}

// strval 数值转字符串
func strval(v any) string {
	switch t := v.(type) {
	case int64:
		return strconv.FormatInt(t, 10)
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	default:
		return strings.TrimSpace(toStr(v))
	}
}

// toStr 尝试调用 String 方法转为字符串
func toStr(v any) string {
	if s, ok := v.(interface{ String() string }); ok {
		return s.String()
	}
	return ""
}
