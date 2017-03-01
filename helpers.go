package melkor

import "github.com/aws/aws-sdk-go/aws"

// ModifyTags takes a set of tags, and modifies them to be more easily filtered
// Input:
//		[
//			{ "Key": "egg", "Value": "bacon" },
//			{ "Key": "bob", "Value": "hope" },
//		]
// Result:
//		[
//			{ "Key": "egg", "Value": "bacon", "egg": "bacon" },
//			{ "Key": "bob", "Value": "hope", "bob": "hope" },
//		]
func ModifyTags(tags interface{}) {
	for _, t0 := range tags.([]interface{}) {
		tag := t0.(map[string]interface{})
		key := aws.StringValue(tag["Key"].(*string))
		value := aws.StringValue(tag["Value"].(*string))
		tag[key] = value
	}
}
