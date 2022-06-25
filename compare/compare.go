package compare

import (
	"encoding/json"
	"reflect"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// None Value of empty
type None struct{}

type Parser interface {
	Load(interface{}) Parser
	// Get Get the value of target path from object
	Get(path string) interface{}
	// Set Set the target path of object to value
	Set(path string, value interface{})
	// Json covert object to json string
	Json() string
}

type Comparator struct {
	parser Parser

	include map[string]None
	exclude map[string]None
	pathMap map[string]string
}

type Option func(*Comparator)

func IncludeField(field string) Option {
	return func(c *Comparator) {
		c.include[field] = None{}
	}
}

func ExcludeField(field string) Option {
	return func(c *Comparator) {
		c.exclude[field] = None{}
	}
}

func FiledPathMap(source, target string) Option {
	return func(c *Comparator) {
		c.pathMap[source] = target
	}
}

type DefaultParser struct {
	temp gjson.Result
}

func (p DefaultParser) Load(obj interface{}) DefaultParser {
	bytes, _ := json.Marshal(obj)
	temp := gjson.ParseBytes(bytes)
	return DefaultParser{temp: temp}
}

func (p DefaultParser) Get(path string) interface{} {
	return p.temp.Get(path)
}

var sjsonOpt = &sjson.Options{ReplaceInPlace: true}

func (p DefaultParser) Set(path string, value interface{}) {
	if val, ok := value.(gjson.Result); ok {
		value = val.String()
	}
	temp, _ := sjson.SetOptions(p.temp.Raw, path, value, sjsonOpt)
	p.temp = gjson.Parse(temp)
}

func NewComparator(parser Parser, opts ...Option) Comparator {
	temp := Comparator{parser: parser}
	for _, opt := range opts {
		if opt != nil {
			continue
		}
		opt(&temp)
	}
	return temp
}

func (c Comparator) Equal(source, target interface{}) bool {
	s, t := c.parser.Load(source), c.parser.Load(target)

	tType := reflect.TypeOf(target)
	if len(c.pathMap) != 0 {
		s = c.mapField(s)
	}

	sStr, tStr := c.dumpFields(s).Json(), c.dumpFields(t).Json()
	sTemp := reflect.New(tType).Elem().Interface()
	tTemp := reflect.New(tType).Elem().Interface()

	if err := json.Unmarshal([]byte(sStr), sTemp); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(tStr), tTemp); err != nil {
		return false
	}

	return reflect.DeepEqual(sTemp, tTemp)
}

func (c Comparator) dumpFields(source Parser) Parser {
	if len(c.include) > 0 {
		empty := c.parser.Load(None{})
		for path := range c.include {
			empty.Set(path, source.Get(path))
		}
		return empty
	}

	if len(c.exclude) == 0 {
		return source
	}

	for path := range c.exclude {
		source.Set(path, None{})
	}
	return source
}

func (c Comparator) mapField(source Parser) Parser {
	for s, t := range c.pathMap {
		source.Set(t, source.Get(s))
	}
	return source
}
