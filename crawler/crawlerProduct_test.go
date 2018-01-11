package crawler

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetCategory(t *testing.T) {
	categoryGtZeroUrl := "https://www.amazon.com/Norton-Facsimile-First-Folio-Shakespeare/dp/0393098435/ref=cm_cr_srp_d_product_top?ie=UTF8"
	categoryEqulZeroUrl := "https://www.amazon.com/Anchorman-Legend-Burgundy-Will-Ferrell/dp/B008LXZEYU/ref=cm_cr_srp_d_product_top?ie=UTF8"
	Convey("test get category", t, func() {

		Convey("When get product info for amazon", func() {
			c, err := getProductCategory(categoryGtZeroUrl)
			So(err, ShouldBeNil)
			So(len(c), ShouldBeGreaterThan, 0)

			c, err = getProductCategory(categoryEqulZeroUrl)
			So(err, ShouldBeNil)
			So(len(c), ShouldEqual, 0)
		})
	})
}
