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
	JAPAN: "session-id=355-9214000-6210429; ubid-acbjp=357-4101823-5570864; x-acbjp=\"fDZLPIY3JAOtazLVEfgaMRQxEg0xS7HOQGzRmchVcVd5FfAJ?MTMGJxmc0KUW3Hx\"; at-acbjp=Atza|IwEBIGVw1hULIcr8B70GBz-gWUVP4Sprr8CvmAY5GP6VSDKWw6IHB6Jp3PyzT2Gv2WAgkSDSznnwT5ZRHGoWNG3GrCh4e48mid5FVJLOBVs3H1kQfU0FxFdHlJPqJ-HbMWvlxZE_-1NGCZwFTIeQdI09BP4snzrQSnhqWtKUoPr_A4rE5T_IZS3CB3MnbM_NmkoDuVEEsxdrJESY6qRV1cUqj-f6U-T0OETpgwfB6YXn2IaLUmhhwr-J9TieubvkWCWHEfLbGYxyMVBvENI6b7uKudYrkXW3zrBhc__hLnqiZl_6H9N2SeOiZP4x4YugOaEHonOkb0zjt_vXr9ICvcj6FtKQIxVWONnnBs9utnn3uHSLl8o0k8rvBRAUurxSpTRfjtG-D8Q-3GN316fZLuf7GVGH; sst-acbjp=Sst1|PQFBLlccABs1Ooy8-wY1U09YCInG7QjspoRI6tKxcCpAVEgO3FfBtSRUoXmN_-di6TYCJPV1kNjQUpLNp7tvwR7QpeseITqnauGoZPF6ju8E6kgYMYUgI33azc5EejJktjZbRTRjz9dQdIl99Rms42iapjvY9tkDxzoBYFFV27zQvNKzqRQVlJigz8MHEHdGBo0hC3dBKeuqlp6ARZPsWzBADILUhqwTKsyLG86j9328LZd0g8rqRS1Cz4JpA8FTmruXxV0Al5qPR6WDBvUCDAwAxw; x-wl-uid=1b/ge7jISAvBGTcaqV4YtpT+sU8RW3r/iEX4R+JuspaxEZYVjQ4RzEWS6joCvh4uBlhg9GKo7ccdpGNhK1A7bH2pAjB77ryDnafA/Bamrpik5BaIcsAYLwEUBv4nhnJWIFI4iWGx7EDc=; lc-acbjp=en_US; session-token=\"OR1sBFquXRTRof269jhozaZUJAPJ6Ye4cmIq/lq3Sj9EwmRUDd7znF9dfLIuzv917PSSIvi2i1Q4onO5VDl2zAO3e/8zr44acVB9Tmilz+47UOP/ovT76Sd6N9jBpJyKd7gm3vTPCDjBFQAsj7sZNses2x2GVc6pFKEjcNmHQvbzTbVeJdEUURTaru9Dqg5rTZpunbhg4PI6of0MArsOdNZMSh/JbEcD/vtWoZJMMjOOwQWcoUgu+3xx3JyM/qK5LNNGicPioWg=\"; session-id-time=2082726001l; csm-hit=s-EM5CNF0D5HF3P758XC19|1516333460437",
	US:    "session-id=133-8500956-2009409; ubid-main=133-6893772-5740726; lc-main=en_US; s_nr=1515462747823-Repeat; s_vnum=1947399226604%26vn%3D2; s_dslv=1515462747825; x-main=\"21E2XUVCUDPsAx9Wc09rmiEsdUGuzzpttTnkp75b4OI9xd4g7a5kpJBLMtq6@Coi\"; at-main=Atza|IwEBIKp7Uiz3lJDJb8cm5_-jDr4D6aNuNNhDFgEhbYpDbipfk9oXWlYQuQIT7M-1bGG58hP-wOTmeQfxPM125tDN9QBsoUHLuhZUNGCLHUs6OHf8vN7ARmnjZT34yti2i91PvTC9cg81vviFaT8SEpk8uc9TIWbgzOfAeA_WqWaW2i_4EsEI-RH7fXrpVsdZGazVlB4HfNInfMB_5GiHOAy3foqJg_7KYGcNVshSotXmMvRKhMph-H_x9dKggBvjjNBY-wpbum-iMYnJBjWk3oPim0I4Pj1jz26pKEdCl7w-xEkyH_ztpxWTZHoDh6rZPNMCTb_anN_uBd8E4I8iB5poiu2sip5jVhrD0b7R71fzTqtYUzYUwqC-02Jn5kz-afLbH_vCT7K0HrymNqLFEl9h7Ge0; sst-main=Sst1|PQFu-vb2AqQdbRCS7iTMaMrZCF5FPHoaCkhCne2n4YKDQ3P8CcG3jpJ1lD0nVApo4gvCG-7cdeTdFhm2jFx3IYqw9dNw_xQqIzJ4BkZ_O9Bgq4---GqQLrLBZOJpyI_qETcd-qU8Z7y9kNTa4tYxzB_xzleV5uy7aGiDmMdMw_4Tp38myVs6jxxx3N1MQZNm17S25-d-4TbdYKC5Vzeeym-IBhwNPh5pPba9Jmrbd94j28gfjpG2lqUDEN6EeNWeUZ5-HRnG7Xmm7U-tFTb9SeNqnA; x-wl-uid=1nqNEvPDY+AGW9Oenk1wtTqUTSVxvIHeL8o96IMjJM6N8ENeCNBI9E+g/9HGn7Fk2iLQ+UjoKB6RqewIkS47wBWh0Y7rYXHmKYMaULZW91+t3ANWuS2hMzfXA/7QIRfuEif/W4zrjYuw=; session-id-time=2082787201l; session-token=SCnwJ6bCuIMnQBvMxMhK/edIB/83b1xOQr691xJ8UW9o7qxzcox7S3gs5rEdfLCvxGJZJCz0U2z0xYy1ejbPnWiKDWzFzJZgSSgtRf85c2o38xYF2CNgw+DdDS8O0gWSGUPb9JciZ4dy6/UZ+iJRpk/CMZkeEzKJmUJqfupJBRaj4oxWHfKePPF47XMf1UnXPLDxmSqjaXKUDj8yQ7JkbSZDTnf793sOavnrJBNIl8FlgsPuwd+J5IVSryA8N9VLgf1Q7qCFIy4gSqBsGmjwVw==; csm-hit=s-M3NTYXHMTSFV593S2M41|1516333570969",
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
