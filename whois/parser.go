//https://help.aliyun.com/knowledge_detail/35772.html
package whois

import (
	"fmt"
	"reflect"
	"strings"
)

type WhoisInfo struct {
	DomainName     string   `whois:"domain name"`
	RegistrarID    string   `whois:"registrar iana id"`
	Status         []string `whois:"domain status,state,status"`
	NameServer     []string `whois:"name server,nserver"`
	UpdatedDate    string   `whois:"updated date,update date"`
	CreationDate   string   `whois:"registration time,creation date"`
	ExpirationDate string   `whois:"expiration date,expiration time,expiry date"`

	WhoisServer     string `whois:"whois server"`
	RegistrantEmail string `whois:"registrant email,registrant contact email"`
	Registrant      string `whois:"registrant name,registrant"`
	//Registrar      RegistrarInfo //注册商
	//Registrant Info `whois:"registrant "` //注册者
	//Administrative Info   `whois:"admin "`      //管理联系人
	//Technical      Info   `whois:"tech "`       //技术联系人
	//Billing        Info   `whois:"billing "`    //付费联系人
	//DNSSEC         string `whois:"dnssec"`
}

//注册商
type RegistrarInfo struct {
	WhoisServer string `whois:"whois server"`
	URL         string `whois:"registrar url,referral url"`
	Name        string `whois:"registrar,registrar"`
	ID          string `whois:"registrar iana id"`
	Email       string `whois:"registrar abuse contact email"`
	Phone       string `whois:"registrar abuse contact phone"`
}

type Info struct {
	ID           string `whois:"id"`
	Name         string `whois:"name"`
	Organization string `whois:"organization"`
	Street       string `whois:"street"`
	City         string `whois:"city"`
	Province     string `whois:"state/province"`
	PostalCode   string `whois:"postal code"`
	Country      string `whois:"country"`
	Phone        string `whois:"phone"`
	PhoneExt     string `whois:"phone ext"`
	Fax          string `whois:"fax"`
	FaxExt       string `whois:"fax ext"`
	Email        string `whois:"email,contact email"`
}

func PrintTags() {
	w := new(WhoisInfo)
	v := reflect.ValueOf(w).Elem()
	t := v.Type()
	for index := 0; index < v.NumField(); index++ {
		vField := v.Field(index)
		tags := strings.Split(t.Field(index).Tag.Get("whois"), ",")
		// name := t.Field(index).Name
		kind := vField.Kind()
		if kind == reflect.Struct {
			Struct := vField
			tStruct := Struct.Type()
			for i := 0; i < Struct.NumField(); i++ {
				// field := Struct.Field(i)
				_tags := strings.Split(tStruct.Field(i).Tag.Get("whois"), ",")
				for _, tag := range tags {
					for _, _tag := range _tags {
						if tag != "" {
							fmt.Println(tag + " " + _tag)
						} else {
							fmt.Println(_tag)
						}

					}
				}
			}
		} else if kind == reflect.Slice {
			for _, tag := range tags {
				fmt.Println(tag)
			}
		} else {
			for _, tag := range tags {
				fmt.Println(tag)
			}
		}
	}
}
func Parse(result string) *WhoisInfo {
	w := new(WhoisInfo)
	v := reflect.ValueOf(w).Elem()
	t := v.Type()
	for index := 0; index < v.NumField(); index++ {
		vField := v.Field(index)
		tags := strings.Split(t.Field(index).Tag.Get("whois"), ",")

		kind := vField.Kind()
		if kind == reflect.Struct {
			Struct := vField
			tStruct := Struct.Type()

			for i := 0; i < Struct.NumField(); i++ {
				_tags := strings.Split(tStruct.Field(i).Tag.Get("whois"), ",")
			NextField:
				for _, tag := range tags {
					for _, _tag := range _tags {
						if tag != "" {
							// fmt.Println(tag + " " + _tag)
							if value, ok := getValue(result, tag+" "+_tag); ok {
								Struct.Field(i).SetString(value)
								break NextField
							}
						} else {
							// fmt.Println(_tag)
							if value, ok := getValue(result, _tag); ok {
								Struct.Field(i).SetString(value)
								break NextField
							}
						}
					}
				}
			}
		} else if kind == reflect.Slice {
			for _, tag := range tags {
				if value, ok := getValueSlice(result, tag); ok {
					vField.Set(reflect.ValueOf(value))
					break
				}
			}
		} else {
			for _, tag := range tags {
				if value, ok := getValue(result, tag); ok {
					vField.SetString(value)
					break
				}
			}
		}
	}
	return w
}

func getValue(result, tag string) (string, bool) {
	key := strings.TrimSpace(tag) + ":"
	start := strings.Index(result, key)
	if start < 0 {
		return "", false
	}
	start += len(key)
	end := strings.Index(result[start:], "\n")
	value := strings.TrimSpace(result[start : start+end])
	return value, true
}

func getValueSlice(result, tag string) (slice []string, ok bool) {
	key := tag + ":"
	for {
		start := strings.Index(result, key)
		if start < 0 {
			break
		}
		ok = true
		start += len(key)
		end := strings.Index(result[start:], "\n")
		value := strings.TrimSpace(result[start : start+end])
		slice = append(slice, value)
		result = result[start+end:]
	}
	return
}
