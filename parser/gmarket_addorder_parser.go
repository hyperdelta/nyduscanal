package parser

import (
	"github.com/buger/jsonparser"
	"strconv"
	"strings"
	"bytes"
)


func GmarketAddOrderParser(data []byte) []byte {
	var result bytes.Buffer

	method, _ := jsonparser.GetString(data, "method")
	if(method != "addOrder") {
		return result.Bytes()
	}
	result.WriteString("{")
	// string으로 한번 하면 앞에 " 이걸 떼서 그런지 바로 []byte로 받는 거랑 값이 다름
	payload,_,_, _ := jsonparser.Get(data, "payload")
	body, _, _, _ := jsonparser.Get(payload, "body")
	// 복수 배송지 걸러내는 로직이 필요함
	addressCount := 0
	jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if(addressCount > 1) {
			return
		}
		addressCount = addressCount + 1 // 복수배송지면 제껴야 함 아에 주지를 않음
		// 메인 주소 2개 가져오기
		deliveryAddr1,_ := jsonparser.GetString(value, "DeliveryAddr1")
		if len(deliveryAddr1) > 0 {
			addrs := strings.Split(deliveryAddr1, " ")
			shippingAddress := addrs[0] + " " + addrs[1]
			result.WriteString("\"ShippingAddress\" : \"" + shippingAddress + "\",")
		}


		// 우편번호 가져오기
		//DeliveryZipCode1
		deliveryZipCode1,_ := jsonparser.GetString(value, "DeliveryZipCode1")
		deliveryZipCode2,_ := jsonparser.GetString(value, "DeliveryZipCode2")
		deliveryZipCode := deliveryZipCode1
		if len(deliveryZipCode2) > 0 {
			deliveryZipCode = deliveryZipCode + "-" + deliveryZipCode2
		}
		result.WriteString("\"DeliveryZipCode\" : \"" + deliveryZipCode + "\",")



		// smilebox 정보 가져오기
		isSmilebox,_ := jsonparser.GetString(value, "IsSmilebox")
		if(len(isSmilebox) < 1) {
			isSmilebox = "false"
		}
		result.WriteString("\"IsSmilebox\" : \"" + isSmilebox + "\",")


		smileboxBranchNo,_ := jsonparser.GetString(value, "SmileboxBranchNo")
		result.WriteString("\"SmileboxBranchNo\" : \"" + smileboxBranchNo + "\",")

	}, "shippingAddressList")


	paymentData, _, _, _ := jsonparser.Get([]byte(body), "paymentData")

	paymentBase, _, _, _ := jsonparser.Get([]byte(paymentData), "PaymentBase")

	// 결제수단
	mediumMethodCode, _ := jsonparser.GetString(paymentBase, "MediumMethodCode")
	result.WriteString("\"MediumMethodCode\" : \"" + mediumMethodCode + "\",")


	smallMethodCode, _ := jsonparser.GetString(paymentBase, "SmallMethodCode")
	result.WriteString("\"SmallMethodCode\" : \"" + smallMethodCode + "\",")

	// 이 부분을 Float로 하면 다시 string으로 바꿔야 하는데 parsing 자체가 안됨 ㅋ
	authPayPrice, _ := jsonparser.GetFloat(paymentBase, "AuthPayPrice")
	result.WriteString("\"AuthPayPrice\" : " + strconv.FormatFloat(authPayPrice, 'f', 0, 64))


	result.WriteString("}")

	if(addressCount > 1) {
		result.Reset()
	}

	return result.Bytes()
}
