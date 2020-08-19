

<p align="center" >Fastjson-based secondary packaging</a></p>

It's easy to construct and parse JSON [How to use](#How-to-use)

Getting Started
===============

## Installing

To start using GJSON, install Go and run `go get`:

```
$ go get -u github.com/yushaona/gjson
```

## How to use
### Construct JSON

- Test code
```
package main

import (
	"fmt"
	
	"github.com/yushaona/gjson"
)
func main() {
	var main gjson.GJSON
	main.SetInt("BigIntSupported", 9223372036854775807)
	main.SetString("date", "20180322")
	main.SetString("city", "北京")

	var data gjson.GJSON
	data.SetFloat64("pm25", 73.5)
	data.SetString("ganmao", "极少数敏感人群应减少户外活动")
	var yesterday gjson.GJSON
	yesterday.SetString("date", "21日星期三")
	yesterday.SetInt("aqi", 85)
	data.SetObject("yesterday", yesterday)
	main.SetObject("data", data)

	var forecast gjson.GJSON
	item1 := forecast.AddItem()
	item1.SetInt("aqi", 98)
	item1.SetString("notice", "愿你拥有比阳光明媚的心情")
	item2 := forecast.AddItem()
	item2.SetInt("aqi", 118)
	item2.SetString("notice", "阴晴之间，谨防紫外线侵扰")
	main.SetArray("forecast", forecast)
	fmt.Println(main.ToString())
}
```
- Output
```
{
    "BigIntSupported": 9223372036854775807,
    "date": "20180322",
    "city": "北京",
    "data": {
        "pm25": 73.5,
        "ganmao": "极少数敏感人群应减少户外活动",
        "yesterday": {
            "date": "21日星期三",
            "aqi": 85
        }
    },
    "forecast": [
        {
            "aqi": 98,
            "notice": "愿你拥有比阳光明媚的心情"
        },
        {
            "aqi": 118,
            "notice": "阴晴之间，谨防紫外线侵扰"
        }
    ]
}
```

### parse JSON
- Test code
```
package main

import (
	"fmt"
	
	"github.com/yushaona/gjson"
)
func main() {
	var main gjson.GJSON
	err := main.Load(`{"BigIntSupported":9223372036854775807,"date":"20180322","city":"北京","data":{"pm25":73.5,"ganmao":"极少数敏感人群应减少户外活动","yesterday":{"date":"21日星期三","aqi":85}},"forecast":[{"aqi":98,"notice":"愿你拥有比阳光明媚的心情"},{"aqi":118,"notice":"阴晴之间，谨防紫外线侵扰"}]}`)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("%s:%d\n", "BigIntSupported", main.GetInt64("BigIntSupported"))
		fmt.Printf("%s:%s\n", "date", main.GetString("date"))
		fmt.Printf("%s:%d\n", "date", main.GetInt64("date"))
		fmt.Printf("%s:%s\n", "city", main.GetString("city"))

		data := main.GetObject("data")
		fmt.Printf("%s:%f\n", "pm25", data.GetFloat64("pm25"))
		fmt.Printf("%s:%s\n", "ganmao", data.GetString("ganmao"))
		yesterday := data.GetObject("yesterday")
		fmt.Printf("%s:%s\n", "date", yesterday.GetString("date"))
		fmt.Printf("%s:%d\n", "aqi", yesterday.GetInt("aqi"))

		forecast := main.GetArray("forecast")
		num := forecast.ItemCount()

		for i := 0; i < num; i++ {
			item := forecast.Item(i)
			fmt.Printf("%s:%d\n", "aqi", item.GetInt("aqi"))
			fmt.Printf("%s:%s\n", "notice", item.GetString("notice"))
		}

	}
}

```
- Output 
```
BigIntSupported:9223372036854775807
date:20180322
date:20180322
city:北京
pm25:73.500000
ganmao:极少数敏感人群应减少户外活动
date:21日星期三
aqi:85
aqi:98
notice:愿你拥有比阳光明媚的心情
aqi:118
notice:阴晴之间，谨防紫外线侵扰
```