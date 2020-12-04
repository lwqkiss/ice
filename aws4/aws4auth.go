package aws4

import (
	"encoding/base64"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Credentials stores the information necessary to authorize with AWS and it
// is from this information that requests are signed.
type Credentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SecurityToken   string `json:"Token"`
	Expiration      time.Time
}

// Sign signs a request bound for AWS. It automatically chooses the best
// authentication scheme based on the service the request is going to.
func Sign(request *http.Request, credentials ...Credentials) *http.Request {
	service, _ := serviceAndRegion(request.URL.Host)
	signVersion := awsSignVersion[service]

	switch signVersion {
	case 2:
		// return Sign2(request, credentials...)
	case 3:
		// return Sign3(request, credentials...)
	case 4:
		return Sign4(request, credentials...)
	case -1:
		return SignS3(request, credentials...)
	}

	return nil
}

// Sign4 signs a request with Signed Signature Version 4.
func Sign4(request *http.Request, credentials ...Credentials) *http.Request {
	keys := chooseKeys(credentials)

	// Add the X-Amz-Security-Token header when using STS
	if keys.SecurityToken != "" {
		request.Header.Set("X-Amz-Security-Token", keys.SecurityToken)
	}

	prepareRequestV4(request)
	meta := new(metadata)

	// Task 1
	hashedCanonReq := hashedCanonicalRequestV4(request, meta)

	// Task 2
	stringToSign := stringToSignV4(request, hashedCanonReq, meta)

	// Task 3
	signingKey := signingKeyV4(keys.SecretAccessKey, meta.date, meta.region, meta.service)
	signature := signatureV4(signingKey, stringToSign)

	request.Header.Set("Authorization", buildAuthHeaderV4(signature, meta, keys))

	return request
}

// SignS3 signs a request bound for Amazon S3 using their custom
// HTTP authentication scheme.
func SignS3(request *http.Request, credentials ...Credentials) *http.Request {
	keys := chooseKeys(credentials)

	// Add the X-Amz-Security-Token header when using STS
	if keys.SecurityToken != "" {
		request.Header.Set("X-Amz-Security-Token", keys.SecurityToken)
	}

	prepareRequestS3(request)

	stringToSign := stringToSignS3(request)
	signature := signatureS3(stringToSign, keys)

	authHeader := "AWS " + keys.AccessKeyID + ":" + signature
	request.Header.Set("Authorization", authHeader)

	return request
}

// SignS3Url signs a GET request for a resource on Amazon S3 by appending
// query string parameters containing credentials and signature. You must
// specify an expiration date for these signed requests. After that date,
// a request signed with this method will be rejected by S3.
func SignS3Url(request *http.Request, expire time.Time, credentials ...Credentials) *http.Request {
	keys := chooseKeys(credentials)

	stringToSign := stringToSignS3Url("GET", expire, request.URL.Path)
	signature := signatureS3(stringToSign, keys)

	query := request.URL.Query()
	query.Set("AWSAccessKeyId", keys.AccessKeyID)
	query.Set("Signature", signature)
	query.Set("Expires", timeToUnixEpochString(expire))
	request.URL.RawQuery = query.Encode()

	return request
}

// expired checks to see if the temporary credentials from an IAM role are
// within 4 minutes of expiration (The IAM documentation says that new keys
// will be provisioned 5 minutes before the old keys expire). Credentials
// that do not have an Expiration cannot expire.
func (this *Credentials) expired() bool {
	if this.Expiration.IsZero() {
		// Credentials with no expiration can't expire
		return false
	}
	expireTime := this.Expiration.Add(-4 * time.Minute)
	// if t - 4 mins is before now, true
	if expireTime.Before(time.Now()) {
		return true
	} else {
		return false
	}
}

type metadata struct {
	algorithm       string
	credentialScope string
	signedHeaders   string
	date            string
	region          string
	service         string
}

const (
	envAccessKey       = "AWS_ACCESS_KEY"
	envAccessKeyID     = "AWS_ACCESS_KEY_ID"
	envSecretKey       = "AWS_SECRET_KEY"
	envSecretAccessKey = "AWS_SECRET_ACCESS_KEY"
	envSecurityToken   = "AWS_SECURITY_TOKEN"
)

var (
	awsSignVersion = map[string]int{
		"autoscaling":          4,
		"ce":                   4,
		"cloudfront":           4,
		"cloudformation":       4,
		"cloudsearch":          4,
		"monitoring":           4,
		"dynamodb":             4,
		"ec2":                  4,
		"elasticmapreduce":     4,
		"elastictranscoder":    4,
		"elasticache":          4,
		"es":                   4,
		"glacier":              4,
		"kinesis":              4,
		"redshift":             4,
		"rds":                  4,
		"sdb":                  2,
		"sns":                  4,
		"sqs":                  4,
		"s3":                   4,
		"elasticbeanstalk":     4,
		"importexport":         4,
		"iam":                  4,
		"route53":              4,
		"elasticloadbalancing": 4,
		"email":                4,
	}
)

func signatureS3(stringToSign string, keys Credentials) string {
	hashed := hmacSHA1([]byte(keys.SecretAccessKey), stringToSign)
	return base64.StdEncoding.EncodeToString(hashed)
}

func stringToSignS3(request *http.Request) string {
	str := request.Method + "\n"

	if request.Header.Get("Content-Md5") != "" {
		str += request.Header.Get("Content-Md5")
	} else {
		body := readAndReplaceBody(request)
		if len(body) > 0 {
			str += hashMD5(body)
		}
	}
	str += "\n"

	str += request.Header.Get("Content-Type") + "\n"

	if request.Header.Get("Date") != "" {
		str += request.Header.Get("Date")
	} else {
		str += timestampS3()
	}

	str += "\n"

	canonicalHeaders := canonicalAmzHeadersS3(request)
	if canonicalHeaders != "" {
		str += canonicalHeaders
	}

	str += canonicalResourceS3(request)

	return str
}

func stringToSignS3Url(method string, expire time.Time, path string) string {
	return method + "\n\n\n" + timeToUnixEpochString(expire) + "\n" + path
}

func timeToUnixEpochString(t time.Time) string {
	return strconv.FormatInt(t.Unix(), 10)
}

func canonicalAmzHeadersS3(request *http.Request) string {
	var headers []string

	for header := range request.Header {
		standardized := strings.ToLower(strings.TrimSpace(header))
		if strings.HasPrefix(standardized, "x-amz") {
			headers = append(headers, standardized)
		}
	}

	sort.Strings(headers)

	for i, header := range headers {
		headers[i] = header + ":" + strings.Replace(request.Header.Get(header), "\n", " ", -1)
	}

	if len(headers) > 0 {
		return strings.Join(headers, "\n") + "\n"
	} else {
		return ""
	}
}

func canonicalResourceS3(request *http.Request) string {
	res := ""

	if isS3VirtualHostedStyle(request) {
		bucketname := strings.Split(request.Host, ".")[0]
		res += "/" + bucketname
	}

	res += request.URL.Path

	for _, subres := range strings.Split(subresourcesS3, ",") {
		if strings.HasPrefix(request.URL.RawQuery, subres) {
			res += "?" + subres
		}
	}

	return res
}

func prepareRequestS3(request *http.Request) *http.Request {
	request.Header.Set("Date", timestampS3())
	if request.URL.Path == "" {
		request.URL.Path += "/"
	}
	return request
}

// Info: http://docs.aws.amazon.com/AmazonS3/latest/dev/VirtualHosting.html
func isS3VirtualHostedStyle(request *http.Request) bool {
	service, _ := serviceAndRegion(request.Host)
	return service == "s3" && strings.Count(request.Host, ".") == 3
}

func timestampS3() string {
	return now().Format(timeFormatS3)
}

const (
	timeFormatS3   = time.RFC1123Z
	subresourcesS3 = "acl,lifecycle,location,logging,notification,partNumber,policy,requestPayment,torrent,uploadId,uploads,versionId,versioning,versions,website"
)
