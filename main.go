package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type InstanceInfo struct {
	InstanceID       string `json:"instanceid"`
	PrivateIpAddress string `json:"privateipaddress"`
	PublicIpAddress  string `json:"publicipaddress"`
}

type AWSResponse struct {
	Reservations []struct {
		Instances []struct {
			InstanceInfo
		} `json:"instances"`
	} `json:"reservations"`
}

type AWSConfig struct {
	client *ec2.Client
}

type FlagsConfig struct {
	Profile string
	Region  string
	IP_Type string
}

func flagSetup() *FlagsConfig {
	profile := flag.String("profile", "default", "the profile to use")
	region := flag.String("region", "us-east-1", "region to use")
	ip_type := flag.String("ip_type", "public", "private or public ip to get from aws")
	flag.Parse()
	return &FlagsConfig{
		Profile: *profile,
		Region:  *region,
		IP_Type: *ip_type,
	}
}

func writeLogs(s InstanceInfo, required string) error {
	var IPAddress string
	file, err := os.OpenFile("generatedLogs", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		return err
	}
	switch strings.ToLower(required) {
	case "private":
		IPAddress = s.PrivateIpAddress
	case "public":
		IPAddress = s.PublicIpAddress
	default:
		fmt.Println("Only private ou public ip addresses")
	}
	if IPAddress != "" {
		_, err = file.WriteString(IPAddress + "\n")
		if err != nil {
			return err
		}
		fmt.Printf("Instância: %s\nIP: %s\n", s.InstanceID, IPAddress)
		defer file.Close()
	}
	return nil
}

func loadConfig(region, profile string) (*AWSConfig, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region), config.WithSharedConfigProfile(profile))
	if err != nil {
		return &AWSConfig{}, err
	}
	client := ec2.NewFromConfig(cfg)
	return &AWSConfig{
		client: client,
	}, nil
}

func listAll(client *ec2.Client) (AWSResponse, error) {

	input := &ec2.DescribeInstancesInput{}

	resp, err := client.DescribeInstances(context.Background(), input)
	if err != nil {
		return AWSResponse{}, err
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return AWSResponse{}, fmt.Errorf("Erro ao converter json para slice de byte")
	}

	var data AWSResponse
	if err = json.Unmarshal(jsonResp, &data); err != nil {
		return AWSResponse{}, fmt.Errorf("Erro ao injetar o json na struct")
	}
	return data, nil
}

func main() {
	c := flagSetup()
	AWSInit, err := loadConfig(c.Region, c.Profile)
	if err != nil {
		log.Fatal(err)
	}
	list, err := listAll(AWSInit.client)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range list.Reservations {
		for _, instanceInfo := range v.Instances {
			err = writeLogs(InstanceInfo{instanceInfo.InstanceID, instanceInfo.PrivateIpAddress, instanceInfo.PublicIpAddress}, c.IP_Type)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
