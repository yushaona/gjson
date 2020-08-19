package gjson

import (
	"testing"
)

/*{
    "BigIntSupported": 429496728,
    "date": "20180322",
    "city": "北京",
    "data": {
        "pm25": 73,
        "ganmao": "极少数敏感人群应减少户外活动",
        "yesterday": {
            "date": "21日星期三",
            "aqi": 85
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
}
*/
//TestConstructJSON  create JSON
func TestConstructJSON(t *testing.T) {
	var main GJSON
	main.SetInt("BigIntSupported", 9223372036854775807)
	main.SetString("date", "20180322")
	main.SetString("city", "北京")

	var data GJSON
	data.SetFloat64("pm25", 73.5)
	data.SetString("ganmao", "极少数敏感人群应减少户外活动")
	var yesterday GJSON
	yesterday.SetString("date", "21日星期三")
	yesterday.SetInt("aqi", 85)
	data.SetObject("yesterday", yesterday)
	main.SetObject("data", data)

	var forecast GJSON
	item1 := forecast.AddItem()
	item1.SetInt("aqi", 98)
	item1.SetString("notice", "愿你拥有比阳光明媚的心情")
	item2 := forecast.AddItem()
	item2.SetInt("aqi", 118)
	item2.SetString("notice", "阴晴之间，谨防紫外线侵扰")
	main.SetArray("forecast", forecast)
	t.Log(main.ToString())

}

func TestParseJSON(t *testing.T) {

	var main GJSON
	err := main.Load(`{"BigIntSupported":9223372036854775807,"date":"20180322","city":"北京","data":{"pm25":73.5,"ganmao":"极少数敏感人群应减少户外活动","yesterday":{"date":"21日星期三","aqi":85}},"forecast":[{"aqi":98,"notice":"愿你拥有比阳光明媚的心情"},{"aqi":118,"notice":"阴晴之间，谨防紫外线侵扰"}]}`)
	if err != nil {
		t.Fatal(err.Error())
	} else {
		t.Logf("%s : %d ", "BigIntSupported", main.GetInt64("BigIntSupported"))
		t.Logf("%s : %s ", "date", main.GetString("date"))
		t.Logf("%s : %d ", "date", main.GetInt64("date"))
		t.Logf("%s : %s ", "city", main.GetString("city"))

		data := main.GetObject("data")
		t.Logf("%s : %f", "pm25", data.GetFloat64("pm25"))
		t.Logf("%s : %s", "ganmao", data.GetString("ganmao"))
		yesterday := data.GetObject("yesterday")
		t.Logf("%s : %s", "date", yesterday.GetString("date"))
		t.Logf("%s : %d", "aqi", yesterday.GetInt("aqi"))

		forecast := main.GetArray("forecast")
		num := forecast.ItemCount()

		for i := 0; i < num; i++ {
			item := forecast.Item(i)
			t.Logf("%s : %d", "aqi", item.GetInt("aqi"))
			t.Logf("%s : %s", "notice", item.GetString("notice"))
		}

	}
}

func TestExist(t *testing.T) {
	var main GJSON
	err := main.Load(`{"BigIntSupported":9223372036854775807,"date":"20180322","city":"北京","data":{"pm25":73.5,"ganmao":"极少数敏感人群应减少户外活动","yesterday":{"date":"21日星期三","aqi":85}},"forecast":[{"aqi":98,"notice":"愿你拥有比阳光明媚的心情"},{"aqi":118,"notice":"阴晴之间，谨防紫外线侵扰"}]}`)
	if err != nil {
		t.Fatal(err.Error())
	} else {
		t.Log(main.IsExist("BigIntSupported"))
		t.Log(main.IsExist("BigInt"))
		t.Log(main.IsExist("data"))
		t.Log(main.IsExist("forecast"))
	}
}
