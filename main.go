package main

import (
	"github.com/sagernet/sing-box/common/geosite"
	"github.com/sagernet/sing/common"
	"github.com/v2fly/v2ray-core/v5/app/router/routercommon"
	"google.golang.org/protobuf/proto"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func parse(vGeositeData []byte) (map[string][]geosite.Item, error) {
	vGeositeList := routercommon.GeoSiteList{}
	err := proto.Unmarshal(vGeositeData, &vGeositeList)
	if err != nil {
		return nil, err
	}
	domainMap := make(map[string][]geosite.Item)
	for _, vGeositeEntry := range vGeositeList.Entry {
		domains := make([]geosite.Item, 0, len(vGeositeEntry.Domain)*2)
		for _, domain := range vGeositeEntry.Domain {
			switch domain.Type {
			case routercommon.Domain_Plain:
				domains = append(domains, geosite.Item{
					Type:  geosite.RuleTypeDomainKeyword,
					Value: domain.Value,
				})
			case routercommon.Domain_Regex:
				domains = append(domains, geosite.Item{
					Type:  geosite.RuleTypeDomainRegex,
					Value: domain.Value,
				})
			case routercommon.Domain_RootDomain:
				if strings.Contains(domain.Value, ".") {
					domains = append(domains, geosite.Item{
						Type:  geosite.RuleTypeDomain,
						Value: domain.Value,
					})
				}
				domains = append(domains, geosite.Item{
					Type:  geosite.RuleTypeDomainSuffix,
					Value: "." + domain.Value,
				})
			case routercommon.Domain_Full:
				domains = append(domains, geosite.Item{
					Type:  geosite.RuleTypeDomain,
					Value: domain.Value,
				})
			}
		}
		domainMap[strings.ToLower(vGeositeEntry.CountryCode)] = common.Uniq(domains)
	}
	return domainMap, nil
}

func generate(datFilePath string, output string) error {
	datFile, err := os.Open(datFilePath)
	if err != nil {
		return err
	}
	defer datFile.Close()
	if err != nil {
		return err
	}
	vData, err := io.ReadAll(datFile)
	if err != nil {
		return err
	}
	outputFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	domainMap, err := parse(vData)
	if err != nil {
		return err
	}
	outputPath, _ := filepath.Abs(output)
	os.Stderr.WriteString("write " + outputPath + "\n")
	return geosite.Write(outputFile, domainMap)
}

func main() {
	// You need to put the newest geosite.dat into the root directory
	generate("geosite.dat", "geosite.db")
}
