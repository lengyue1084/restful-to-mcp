package tool

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flow-bridge-mcp/internal/mcp/config"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func Capitalize(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// FileExists 检查文件是否存在
func FileExists(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false, err
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// 去重函数
func Unique(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// PathExists 返回true则不存在，返回err具体分析
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// RemoveStringLastChar 移除字符串的最后一个字符
func RemoveStringLastChar(str string) string {
	if len(str) > 0 {
		return str[:len(str)-1]
	}
	return str
}

// Copy 从一个结构体复制到另一个结构体
// 注意：只会复制可JSON序列化的字段，私有字段不会被复制
func Copy(to, from interface{}) error {
	b, err := json.Marshal(from)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, to)
	if err != nil {
		return err
	}
	return nil
}

var (
	// 预编译正则表达式，匹配 {{content}} 格式
	placeholderRegex = regexp.MustCompile(`\{\{\s*([^{}]+)\s*}}`)
	// 预编译正则表达式，匹配 {{.Args.param}} 格式
	argsPlaceholderRegex = regexp.MustCompile(`\{\{\s*\.Args\.\s*([^{}]+)\s*}}`)
)

// ConvertPathToArgsFormat 将路径中的 {{param}} 转换为 {{.Args.param}} 格式
func ConvertPathToArgsFormat(path string) string {
	// 简单替换版本
	return placeholderRegex.ReplaceAllString(path, "{{.Args.$1}}")
}

// ConvertPathToArgsFormatV2 或者使用更复杂的处理版本
func ConvertPathToArgsFormatV2(path string) string {
	return placeholderRegex.ReplaceAllStringFunc(path, func(matched string) string {
		matches := placeholderRegex.FindStringSubmatch(matched)
		if len(matches) < 2 {
			return matched
		}
		content := strings.TrimSpace(matches[1])
		// 可以在这里添加额外的逻辑处理
		return fmt.Sprintf("{{.Args.%s}}", content)
	})
}

// ConvertArgsToPathFormat 将路径中的 {{.Args.param}} 转换为 {{param}} 格式（简单版本）
func ConvertArgsToPathFormat(path string) string {
	// 简单替换版本
	return argsPlaceholderRegex.ReplaceAllString(path, "{{$1}}")
}

// ConvertArgsToPathFormatV2 逆向转换的复杂处理版本
func ConvertArgsToPathFormatV2(path string) string {
	return argsPlaceholderRegex.ReplaceAllStringFunc(path, func(matched string) string {
		matches := argsPlaceholderRegex.FindStringSubmatch(matched)
		if len(matches) < 2 {
			return matched
		}
		content := strings.TrimSpace(matches[1])
		// 可以在这里添加额外的逻辑处理，比如参数验证、格式化等
		return fmt.Sprintf("{{%s}}", content)
	})
}

func CopyStruct(src, dst interface{}) {
	srcVal := reflect.ValueOf(src).Elem()
	dstVal := reflect.ValueOf(dst).Elem()
	for i := 0; i < srcVal.NumField(); i++ {
		value := srcVal.Field(i)
		name := srcVal.Type().Field(i).Name
		dstValueName := dstVal.FieldByName(name)
		if dstValueName.IsValid() == false {
			continue
		}
		dstValueName.Set(value) //这里默认共同成员的类型一样，否则这个地方可能导致 panic，需要简单修改一下。
	}
}

func MapToJson(param map[string]interface{}) (string, error) {
	dataType, err := json.Marshal(param)
	if err != nil {
		return "", err
	}
	dataString := string(dataType)
	return dataString, nil
}

func JsonToMap(str string) (map[string]interface{}, error) {
	var tempMap map[string]interface{}
	err := json.Unmarshal([]byte(str), &tempMap)
	if err != nil {
		return nil, err
	}
	return tempMap, nil
}

// MD5 32位小写
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	h.Write([]byte("pdCat"))
	return hex.EncodeToString(h.Sum(nil))
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}

func WithPrefix(str, prefix string) string {
	if !strings.Contains(str, prefix) {
		str = prefix + str
	}
	return str
}

func RemovePrefix(str, prefix string) string {
	l := len(prefix)
	return string([]byte(str)[l:])
}

// DifferenceSet sli2中的包含的sli1
func DifferenceSet(sli1 []int64, sli2 []int64) []int64 {
	if sli1 == nil || len(sli1) == 0 {
		return sli2
	}
	sli3 := append(sli1, sli2...)
	var DiffIdsMap = make(map[int64]int64)
	for _, val := range sli3 {
		DiffIdsMap[val] = val
	}
	for key, val := range DiffIdsMap {
		for _, v := range sli1 {
			if val == v {
				delete(DiffIdsMap, key)
				break
			}
		}
	}
	var ids = make([]int64, 0, len(DiffIdsMap))
	for _, val := range DiffIdsMap {
		ids = append(ids, val)
	}
	return ids
}

// FileNameByUUid 随机生成文件名
func FileNameByUUid() (filename string) {
	filename = strconv.FormatInt(time.Now().Unix(), 10) + RandString() + uuid.NewString()
	return
}

// StringNameByUUidWithDate 随机生成文件名
func StringNameByUUidWithDate() (filename string) {
	filename = CurrentDateToString("YmdHis") + RandString() + uuid.NewString()
	return
}

// NewUUID 实现一个uuid的方法
func NewUUID() string {
	return uuid.NewString()
}

// RandInt 随机int
func RandInt() int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(9999)
}

// RandString 随机string
func RandString() string {
	return strconv.Itoa(RandInt())
}

// RandStringByLen  随机指定长度的string
func RandStringByLen(n int) string {
	return fmt.Sprintf("%06d", RandIntByLen(n))
}

// RandIntByLen 获取指定位数的数字int随机数
func RandIntByLen(n int) int {
	var m = 1
	for n > 0 {
		m = m * 10
		n--
	}
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(m - 1)
}

// RandStringWithLowercaseAndDigits 生成包含小写字母和数字的随机字符串
func RandStringWithLowercaseAndDigits(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	// 将随机数生成器提升为全局变量或单例，避免重复创建
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	if n <= 0 {
		return ""
	}
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func FirstNonEmpty(str1, str2 string) string {
	if str1 != "" {
		return str1
	}
	if str2 != "" {
		return str2
	}
	return ""
}

// UniqSlice merges multiple slices with removing duplicates the elements.
//
// Play: https://go.dev/play/p/clvg0gFoBQs
func UniqSlice[T comparable](ss ...[]T) []T {
	if len(ss) == 0 {
		return []T{}
	}
	size := 0
	for _, v := range ss {
		size += len(v)
	}

	if size == 0 {
		return []T{}
	}

	var (
		res   = make([]T, 0, size)
		exist = make(map[T]struct{}, size)
	)
	for _, s := range ss {
		for _, v := range s {
			if _, ok := exist[v]; ok {
				continue
			}
			exist[v] = struct{}{}
			res = append(res, v)
		}
	}
	return res
}

// Union returns the union of multiple slices.
func Union[T comparable](ss ...[]T) []T {
	return UniqSlice[T](ss...)
}

func GetBool(m map[string]any, key string, defaultValue bool) bool {
	if m == nil {
		return defaultValue
	}
	v, ok := m[key]
	if !ok {
		return defaultValue
	}
	b, ok := v.(bool)
	if !ok {
		return defaultValue
	}
	return b
}

func GetString(m map[string]any, key string, defaultValue string) string {
	if m == nil {
		return defaultValue
	}
	v, ok := m[key]
	if !ok {
		return defaultValue
	}
	s, ok := v.(string)
	if !ok {
		return defaultValue
	}
	return s
}
func IsJson(data []byte) bool {
	data = bytes.TrimSpace(data)
	if len(data) <= 0 {
		return false
	}
	return json.Valid(data)
}

// CleanBase64String 清理Base64字符串
func CleanBase64String(content string) string {
	// 移除空白字符
	content = strings.ReplaceAll(content, " ", "")
	content = strings.ReplaceAll(content, "\n", "")
	content = strings.ReplaceAll(content, "\r", "")
	content = strings.ReplaceAll(content, "\t", "")
	return content
}
func ValidateBase64String(content string) error {
	// Base64字符串长度必须是4的倍数
	if len(content)%4 != 0 {
		// 尝试添加填充
		padding := 4 - len(content)%4
		if padding != 4 {
			for i := 0; i < padding; i++ {
				content += "="
			}
		}
	}

	// 检查是否只包含合法的Base64字符
	validBase64Regex := regexp.MustCompile(`^[A-Za-z0-9+/]*={0,2}$`)
	if !validBase64Regex.MatchString(content) {
		return fmt.Errorf("包含非法的Base64字符")
	}

	return nil
}

// TryMultipleBase64Decodings 尝试多种Base64解码方式
func TryMultipleBase64Decodings(content string) ([]byte, error) {
	// 1. 标准Base64解码
	if decoded, err := base64.StdEncoding.DecodeString(content); err == nil {
		return decoded, nil
	}

	// 2. URL安全的Base64解码
	if decoded, err := base64.URLEncoding.DecodeString(content); err == nil {
		return decoded, nil
	}

	// 3. Raw标准Base64解码（无填充）
	if decoded, err := base64.RawStdEncoding.DecodeString(content); err == nil {
		return decoded, nil
	}

	// 4. Raw URL安全的Base64解码（无填充）
	if decoded, err := base64.RawURLEncoding.DecodeString(content); err == nil {
		return decoded, nil
	}

	return nil, fmt.Errorf("无法解码Base64字符串")
}
func SplitByMultipleDelimiters(s string, delimiters ...string) []string {
	if len(delimiters) == 0 {
		return []string{s}
	}
	delimiterPattern := "[" + regexp.QuoteMeta(strings.Join(delimiters, "")) + "]"
	re := regexp.MustCompile(delimiterPattern)
	return re.Split(s, -1)
}

func CurrentDateToString(str string) string {
	now := time.Now()
	switch str {
	case "Y-m-d", "Y-M-D", "y-m-d":
		return now.Format("2006-01-02")
	case "Ymd", "YMD", "ymd":
		return now.Format("20060102")
	case "Y-m-d H:i:s":
		return now.Format("2006-01-02 15:04:05")
	case "YmdHis", "ymdhis":
		return now.Format("20060102150405")
	default:
		return now.Format("2006-01-02")
	}
}

// WriteFile 将内容写入指定路径的文件
func WriteFile(path, fileName string, content []byte) (name string, err error) {
	// 参数验证
	if path == "" {
		return "", fmt.Errorf("文件路径不能为空:path:%s", path)
	}
	path = filepath.Join(path, fileName)
	// 获取文件所在目录
	dir := filepath.Dir(path)
	// 确保目录存在
	if err = os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("创建目录失败 %s: %w", dir, err)
	}
	err = os.WriteFile(path, content, 0644)
	if err != nil {
		return "", fmt.Errorf("写入文件 %s 失败: %w", path, err)
	}
	return path, nil
}

func DateToString(t time.Time, str string) string {
	now := t
	switch str {
	case "Y-m-d", "Y-M-D", "y-m-d":
		return now.Format("2006-01-02")
	case "Ymd", "YMD", "ymd":
		return now.Format("20060102")
	case "Y-m-d H:i:s":
		return now.Format("2006-01-02 15:04:05")
	case "YmdHis", "ymdhis":
		return now.Format("20060102150405")
	case "m/d h:i", "m/d H:i":
		return now.Format("01/02 15:04")
	default:
		return now.Format("2006-01-02")
	}
}

func DateStrToTimeStamp(dateStr, str string) (unix int64, err error) {
	cst, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return 0, err
	}
	CSTLayout := "2006-01-02 15:04:05"
	switch str {
	case "Y-m-d", "Y-M-D", "y-m-d":
		CSTLayout = "2006-01-02"
	case "Ymd", "YMD", "ymd":
		CSTLayout = "20060102"
	case "Y-m-d H:i:s":
		CSTLayout = "2006-01-02 15:04:05"
	case "YmdHis", "ymdhis":
		CSTLayout = "20060102150405"
	default:
		CSTLayout = "2006-01-02"
	}
	ts, err := time.ParseInLocation(CSTLayout, dateStr, cst)
	if err != nil {
		return 0, err
	}
	unix = ts.Unix()
	return
}

func ReadTxtFromFile(path string) (content string, err error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		return "", err
	}
	defer f.Close()
	contentByte, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(contentByte), nil
}

// RFC3339ToCSTLayout convert rfc3339 value to china standard time layout
func RFC3339ToCSTLayout(value string) (string, error) {
	cst, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return "", err
	}
	const CSTLayout = "2006-01-02 15:04:05"
	ts, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return "", err
	}
	return ts.In(cst).Format(CSTLayout), nil
}

func RFC3339ToCSTLayoutInt64(value string) (int64, error) {
	cst, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return 0, err
	}
	const CSTLayout = "2006-01-02 15:04:05"
	ts, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return 0, err
	}
	return ts.In(cst).Unix(), nil
}

func DateToCSTLayout(value string) (int64, error) {
	cst, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return 0, err
	}
	const CSTLayout = "2006-01-02 15:04:05"
	ts, err := time.ParseInLocation(CSTLayout, value, cst)
	if err != nil {
		return 0, err
	}
	return ts.In(cst).Unix(), nil
}

func StrDateToCSTLayoutDate(value string) (*time.Time, error) {
	cst, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return nil, err
	}
	const CSTLayout = "2006-01-02 15:04:05"
	ts, err := time.ParseInLocation(CSTLayout, value, cst)
	if err != nil {
		return nil, err
	}
	return &ts, nil
}

func excelDateToDate(excelDate string) time.Time {
	excelTime := time.Date(1899, time.December, 30, 0, 0, 0, 0, time.UTC)
	var days, _ = strconv.Atoi(excelDate)
	return excelTime.Add(time.Second * time.Duration(days*86400))
}

// SliceRemoveSomeSlice 从outIds中移除hasIds切片中包含的id
func SliceRemoveSomeSlice(hasIds, outIds []int64) map[int]int64 {
	IdsMap := make(map[int]int64, len(outIds))
	for k, v := range outIds {
		IdsMap[k] = v
	}
	for _, v := range hasIds {
		for key, value := range IdsMap {
			if v == value {
				delete(IdsMap, key)
			}
		}
	}
	return IdsMap
}

// SliceRemoveSomeSliceBackSlice 从outIds中移除hasIds切片中包含的切片
func SliceRemoveSomeSliceBackSlice(hasIds, outIds []int64) (backIds []int64) {
	IdsMap := make(map[int]int64, len(outIds))
	for k, v := range outIds {
		IdsMap[k] = v
	}
	for _, v := range hasIds {
		for key, value := range IdsMap {
			if v == value {
				delete(IdsMap, key)
			}
		}
	}
	for _, v := range IdsMap {
		backIds = append(backIds, v)
	}
	return backIds
}

// DownLoad 下载图片信息
// 未验证是否是网络地址
func DownLoad(dirBase string, url string, second int) (string, error) {
	fileBase := dirBase
	idx := strings.LastIndex(url, "/")
	if idx < 0 {
		fileBase += "/" + url
	} else {
		fileBase += url[idx+1:]
	}
	if 0 == second {
		second = 60
	}
	timeout := time.Duration(60 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	v, err := client.Get(url)
	if err != nil {
		return url, err
	}
	defer v.Body.Close()
	content, err := ioutil.ReadAll(v.Body)
	if err != nil {
		fmt.Printf("Read http response failed! %v", err)
		return url, err
	}
	err = ioutil.WriteFile(fileBase, content, 0666)
	if err != nil {
		fmt.Printf("Save to file failed! %v", err)
		return url, err
	}
	return url, nil
}

// WeekIntervalTime 获取某周的开始和结束时间,week为0本周,-1上周，1下周以此类推
func WeekIntervalTime(week int) (startTime, endTime string) {
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	//周日做特殊判断 因为time.Monday = 0
	if offset > 0 {
		offset = -6
	}

	year, month, day := now.Date()
	thisWeek := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	startTime = thisWeek.AddDate(0, 0, offset+7*week).Format("2006-01-02") + " 00:00:00"
	endTime = thisWeek.AddDate(0, 0, offset+6+7*week).Format("2006-01-02") + " 23:59:59"

	return startTime, endTime
}

// MonthIntervalTime 获取某月的开始和结束时间mon为0本月,-1上月，1下月以此类推
func MonthIntervalTime(mon int) (startTime, endTime string) {
	year, month, _ := time.Now().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	startTime = thisMonth.AddDate(0, mon, 0).Format("2006-01-02") + " 00:00:00"
	endTime = thisMonth.AddDate(0, mon+1, -1).Format("2006-01-02") + " 23:59:59"
	return startTime, endTime
}

// SecondToHourMinuteSecond 秒转小时分钟秒
func SecondToHourMinuteSecond(second int64) (hmsMap map[string]int) {
	hmsMap["h"] = int(math.Floor(float64(second / 3600)))
	hmsMap["m"] = int(math.Floor(float64((second % 3600) / 60)))
	hmsMap["s"] = int((second % 3600) % 60)
	return
}

// SecondToHourMinuteSecondStr 秒转小时分钟秒
func SecondToHourMinuteSecondStr(second int64) (hmsStr string) {
	hmsMap := make(map[string]int, 3)
	hmsMap["h"] = int(math.Floor(float64(second / 3600)))
	hmsMap["m"] = int(math.Floor(float64((second % 3600) / 60)))
	hmsMap["s"] = int((second % 3600) % 60)

	return strconv.Itoa(hmsMap["h"]) + "小时" + strconv.Itoa(hmsMap["m"]) + "分钟" + strconv.Itoa(hmsMap["s"]) + "秒"
}

// DirExists  判断所给路径文件/文件夹是否存在
func DirExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// AllLetterArr 获取26个字符数组
func AllLetterArr(letter bool) (letterList [26]string) {
	key := 0
	for c := 'A'; c <= 'Z'; c++ {
		if letter != true {
			letterList[key] = strings.ToUpper(string(c))
		} else {
			letterList[key] = string(c)
		}

		key = key + 1
	}
	return
}

// AllLetterSlice 获取26个字符数组
func AllLetterSlice(letter bool) (letterList []string) {
	for c := 'A'; c <= 'Z'; c++ {
		if letter != true {
			letterList = append(letterList, strings.ToUpper(string(c)))
		} else {
			letterList = append(letterList, string(c))
		}
	}
	return
}

// SomeLetterList 获取26个字符数组
func SomeLetterList(letter bool, start, end int) (letterList [26]string) {

	//todo
	return
}

// ExtractArgsFromPath 提取路径中的参数 /{{.Config.url}}/api-backup/user/{{.Args.userid}}/order/{{.Args.orderId}}/list"
func ExtractArgsFromPath(path string) []config.ArgConfig {
	// 使用预编译的正则表达式匹配 {{.Args.param}} 格式
	matches := argsPlaceholderRegex.FindAllStringSubmatch(path, -1)

	var args []config.ArgConfig
	for _, match := range matches {
		if len(match) >= 2 {
			arg := config.ArgConfig{
				Name:        strings.TrimSpace(match[1]),
				Position:    "path",
				Required:    true,
				Type:        "string",
				Description: "",
				Default:     "",
				Items:       config.ItemsConfig{},
				Enum:        nil,
				Explode:     false,
			}
			args = append(args, arg)
		}
	}
	return args
}
