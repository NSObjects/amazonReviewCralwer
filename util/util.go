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
	JAPAN: "session-id=355-9214000-6210429; ubid-acbjp=357-4101823-5570864; x-wl-uid=1b/ge7jISAvBGTcaqV4YtpT+sU8RW3r/iEX4R+JuspaxEZYVjQ4RzEWS6joCvh4uBlhg9GKo7ccdpGNhK1A7bH2pAjB77ryDnafA/Bamrpik5BaIcsAYLwEUBv4nhnJWIFI4iWGx7EDc=; x-acbjp=\"GCa1Q980MxUo?OgOGVs3dgMAahZiB?O2dm@c9u5c3S5mddovZFS7ItWJfkPgD3qn\"; at-acbjp=Atza|IwEBIITvE1kQc1eTQtQCMuyG6LOKfuPFJG4bdaaBSE3sBGLxvdlL3h_0CgB-WUzTTjcA8GBavoiEilPWZu2vOPfU6NhpYGBT_19eoCI2RCO74_hnksNgnMS1gtPbGUe4Q1-H9az6tMHKqD-zDT8NkG5CTSe_eaRBpm_UnbZ9vz7WaKfLlVkiXOQ8gxwSNiRmmWa46H5gb4SaEinQOqDNUyg5FsLmtr_OFOU-ECrLvQY-EPOJt-zw_t_DCR1bPWhIm5zJ8mElNEdWhwFXuUZ-eQggWtI-XsdQU-DQ6204cOk7gRAscLaIeAWvo8aRanOWMdQi1kWPxasYS8SUWJfMkMhYiIpdHHA90FOLEGr-9Q850XEUfYV77CSJnPlWGxMQkRCNVXQG9SRrA3-LFR9qKgrL793z; sst-acbjp=Sst1|PQF1q6HZuVHzlJZB34qqInytCF8daIFcDrDbzY3M08plAVq3-4lCxCTOCOKXFzlIpb5VKC3Fl2BfXN9mvQqM7wDrQn_Jp5jm4K5faXkH_Uzgq0tAyZj7QVzm-E5w1o8FN9MWSXcj68zNim-8dCRKksjkjYHjC0KmyczSnptkwmlIp53WklpHUHxMEDxN-zen4QT0_Pl0XMb3PnzBtLcE6Ofs63vGmC8JWV20rE2qpfgfy2pXDBUcp0ZnDu6vIkaCCa0RxA7c4QBA7zTQOqaYFPaUcw; lc-acbjp=en_US; session-token=\"dak7X9GydtzyUuQJ6LR5hnu1vEs0AlHza5I5NVMI47gTS8klx9D6/HiUoOo9bMflRgms1QCu8UDtdvlLpUeOgSEdwGrJ+IVo+PBWLPj5I5e1GQaO2LCeCwn3FUHavSGwdETaHx5zVc0vWKISMMOJW+/cN0gsG2rJayEYJ32UNpXuHRqEBjkfzjvCxlPAdbHpFm4273J1JEBj4ta0QS/wcmd4T6/mcGq1gJP0wgXGWp+L9U627wbFsPUBVHFXChy5oAPbx/xCdVo=\"; skin=noskin; session-id-time=2082726001l; csm-hit=s-3F5F5SNCFPVYKWJWPQFA|1517018398963",
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
