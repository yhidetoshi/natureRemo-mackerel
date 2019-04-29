package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mackerelio/mackerel-client-go"
)

var (
	url    = "https://api.nature.global/1/devices"
	token  = os.Getenv("REMOTOKEN")
	mkrKey = os.Getenv("MKRKEY")
	client = mackerel.NewClient(mkrKey)
)

const (
	serviceName = "NatureRemo"
	timezone    = "Asia/Tokyo"
	offset      = 9 * 60 * 60
)

type Device struct {
	ID                string       `json:"id"`
	Name              string       `json:"name"`
	TemperatureOffset int32        `json:"temperature_offset"`
	HumidityOffset    int32        `json:"humidity_offset"`
	CreatedAt         string       `json:"created_at"`
	UpdatedAt         string       `json:"updated_at"`
	FirmwareVersion   string       `json:"firmware_version"`
	NewestEvents      NewestEvents `json:"newest_events"`
}

type NewestEvents struct {
	Temperature SensorValue `json:"te"`
	Humidity    SensorValue `json:"hu"`
	Illuminance SensorValue `json:"il"`
}

type SensorValue struct {
	Value     float64 `json:"val"`
	CreatedAt string  `json:"created_at"`
}

func (d *Device) FetchValesFromNatureRemo() (float64, float64, float64) {
	var devices []*Device

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error")
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error request")
	}

	err = json.NewDecoder(resp.Body).Decode(&devices)
	if err != nil {
		fmt.Println("Error decode")
	}

	return devices[0].NewestEvents.Temperature.Value,
		devices[0].NewestEvents.Humidity.Value,
		devices[0].NewestEvents.Illuminance.Value
}

func PostValuesToMackerel(tem float64, hum float64, ill float64, nowTime time.Time) {

	fmt.Println(nowTime)
	// Post Temperature
	err_tem := client.PostServiceMetricValues(serviceName, []*mackerel.MetricValue{
		&mackerel.MetricValue{
			Name:  "Temperature.temperature",
			Time:  nowTime.Unix(),
			Value: tem,
		},
	})
	if err_tem != nil {
		fmt.Println("Error post tem")
	}

	// Post Humidity
	err_hum := client.PostServiceMetricValues(serviceName, []*mackerel.MetricValue{
		&mackerel.MetricValue{
			Name:  "Humidity.humidity",
			Time:  nowTime.Unix(),
			Value: hum,
		},
	})
	if err_hum != nil {
		fmt.Println("Error post hum")
	}

	// Post Illuminance
	err_ill := client.PostServiceMetricValues(serviceName, []*mackerel.MetricValue{
		&mackerel.MetricValue{
			Name:  "Illuminance.illuminance",
			Time:  nowTime.Unix(),
			Value: ill,
		},
	})
	if err_ill != nil {
		fmt.Println("Error post ill")
	}
}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context) {

	jst := time.FixedZone(timezone, offset)
	nowTime := time.Now().In(jst)
	d := &Device{}
	tem, hum, ill := d.FetchValesFromNatureRemo()
	PostValuesToMackerel(tem, hum, ill, nowTime)
}
