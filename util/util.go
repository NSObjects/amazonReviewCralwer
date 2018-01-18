package util

var BaseUrl string
var Cookie string

type Country int

const (
	JAPAN Country = iota
	US
)

var countryURL = map[Country]string{
	JAPAN: "https://www.amazon.co.jp",
	US:    "https://www.amazon.com",
}

var countryCookie = map[Country]string{
	JAPAN: "session-id=355-9214000-6210429; ubid-acbjp=357-4101823-5570864; a-ogbcbff=1; x-acbjp=\"fDZLPIY3JAOtazLVEfgaMRQxEg0xS7HOQGzRmchVcVd5FfAJ?MTMGJxmc0KUW3Hx\"; at-acbjp=Atza|IwEBIGVw1hULIcr8B70GBz-gWUVP4Sprr8CvmAY5GP6VSDKWw6IHB6Jp3PyzT2Gv2WAgkSDSznnwT5ZRHGoWNG3GrCh4e48mid5FVJLOBVs3H1kQfU0FxFdHlJPqJ-HbMWvlxZE_-1NGCZwFTIeQdI09BP4snzrQSnhqWtKUoPr_A4rE5T_IZS3CB3MnbM_NmkoDuVEEsxdrJESY6qRV1cUqj-f6U-T0OETpgwfB6YXn2IaLUmhhwr-J9TieubvkWCWHEfLbGYxyMVBvENI6b7uKudYrkXW3zrBhc__hLnqiZl_6H9N2SeOiZP4x4YugOaEHonOkb0zjt_vXr9ICvcj6FtKQIxVWONnnBs9utnn3uHSLl8o0k8rvBRAUurxSpTRfjtG-D8Q-3GN316fZLuf7GVGH; sess-at-acbjp=\"04qK2Br/anhP7qEoLVyi1klFZW9R0XE9LBcsRZk0Ou4=\"; sst-acbjp=Sst1|PQFBLlccABs1Ooy8-wY1U09YCInG7QjspoRI6tKxcCpAVEgO3FfBtSRUoXmN_-di6TYCJPV1kNjQUpLNp7tvwR7QpeseITqnauGoZPF6ju8E6kgYMYUgI33azc5EejJktjZbRTRjz9dQdIl99Rms42iapjvY9tkDxzoBYFFV27zQvNKzqRQVlJigz8MHEHdGBo0hC3dBKeuqlp6ARZPsWzBADILUhqwTKsyLG86j9328LZd0g8rqRS1Cz4JpA8FTmruXxV0Al5qPR6WDBvUCDAwAxw; lc-acbjp=zh_CN; x-wl-uid=1AqxZ1hl5PEG9QYPsozv77XcES9R61U6sug/+fyUx114GwfPT7n7bHREhU5wVu7D41r2j2WhsSCS1h3MYyHP/fKiYeoDC07WMw3OzkM6ADAf6o7pfuKyOXIpOYVfAPkz1pTLHoVecycI=; session-token=\"TyKu+nLrh0q7/yAuTPIpVCBcIXWqPFXNEA7UcYXx/MdLwXcXfXMBUeZc7nMi3NJPdC2TvdqPu4kRDQkgCoqh/tvl9zFlKKW4sC99QLlu5RuJJSu7WbMSe4J046ZpaztWPAuQSBbVvqykT+YkAL44pz/+UH+DejfvkNwkceRDh2cXL60aWvbjEXxdyKvKJKOxjgJCVhuWQexfsXKkPeXxm6hJOJ/5yvVRQcOacMNqGajujIoOpsw4/a6c8Nia2Gn/7Z40h+wqVY4=\"; session-id-time=2082726001l; csm-hit=s-2HYGX7QQE9R3FXQD3V89|1516181928306",
	US: "session-id=133-8500956-2009409; " +
		"session-id-time=2082787201l;" +
		" ubid-main=133-6893772-5740726;" +
		" s_cc=true; s_vnum=1947399226604&vn=1;" +
		" s_sq=[[B]]; s_nr=1515399305583-New; " +
		"s_dslv=1515399305584; s_ppv=100;" +
		" a-ogbcbff=1; " +
		"session-token=\"cZSlPMRt8u360s/ufDXaODK8+015vU0hYdlcvdE7Q95nAyjN8k+/2+GOwXFaQstzBtvCMEQAqoLREuME6kihWyFmD0WkzCL7UtAx81U4A3xAstdQtNC8HPa1/jSqtd14RY+eSpu695Lv2VHugVAo+n8qJBlOHAnhqDzSsIJxAyUtFmgGFssmVm7kByuuwil2tIr06Dq0CTFGh2MmBUf1eQ6oIhrVFFmKeFeEP09WlWK8HJzazZj5Kbm57t3l1pxiw6+Jhj5KM9tVsGXscLHY/A==\"; x-main=\"ezZ0qDyjzxneRSCR0IKIci1dA?DOIGkpwjYi0?TFgTbo5Dg72O8@dTNP4IZ@n4my\"; at-main=Atza|IwEBIHT7pByiU9me1KMaFTiD0r4OLIbVik0guwYh3mnGxvLelj2UjJDB-thUzAad5hI62iEKmxSGCWc0taQolNVTfzkhb2bxwHd3L6ldO4Acr2wVZ5cF4IyfOrWeiSVN0-rTo7eldGJR3ufwoFp5mspeKzxOFruM1JZx9x68aRnt3HTk6zpERVQHMpiQnktffDPjs3Yb2sFc0V-lwX1BtSbcZ2uWepwuOl7skiwIlCcIPVjaVVoq_6a6_mgvW2YE2BcZ65mUVPWk9QM3fyQcH17h-fTv-BTBonY_avjCCVbWZaVMMl2hmV4uklrdsog5r7zr8O0QVN1cKqVILo26E8WL0eLkWI6fCLatUn3t86XazjpRZux_mqlWHpAEhBP0cHyX9_CywtbqEt5rQjqIQQxSrn-v; sess-at-main=\"GP+xmj/5F+DSiAzermlfqyR5CwvOl2Zrfn7Iwd8uAtc=\"; sst-main=Sst1|PQHODdLpmKLUPsmWTUs5mpH4CFY1XAjRFaDM8Zly3R6fGRG1PXN-RR-BragU1OorONv41QnmEHIgx3WFk9QoEvyf2564ywr4WADsTp50fWbvQEOAQEnPK8JAX5EeIObF3lQCmJSX-1WK30c1Nj7KYwsrwSpubcvEKl6zw2mRke_DWGaNsJW1ImwLQh79V0V_JuK8B8hQx7SlM8d1EY9r1773HqQlurzRE2ZoHYv3RZ33QJVW8ycvqp9zZNzuDN6AFg1rC68jo08dst4_tGZ5CX57fg; lc-main=en_US; x-wl-uid=1bjXOXOzkydlFDdbs+hsDH18XuaxumxxDhvrejiwa7fbc9AC0HCkALwwstLnbInpMjOZwkf6ojugkuvXHP0yKP25KmErlFwTxwEpUXrOsvdu0o6MhOxn+sIEs56r6sucX5GY7CD2/y7w=; csm-hit=s-KYRJW9N1Y95HRP7BHNJ5|1515403111471",
}

func Substr(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < 0 || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}

func SetCountry(c Country) {
	Cookie = countryCookie[c]
	BaseUrl = countryURL[c]
}
