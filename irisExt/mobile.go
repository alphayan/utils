/**
 * <p>Description: (移动接口) </>
 * @author lizhi_duan
 * @date 2018/9/3 16:14
 * @version 1.0
 */
package xcjIris

import (
	"github.com/alphayan/iris/context"
	"strconv"
)

const (
	MASTER_PLATFROM  = "MASTER"
	STUDENT_PLATFORM = "STUDENT"
	PARENT_PLATFORM  = "PARENT"
)

/**
 * <p>Description: (验证是否是移动端登录) </p>
 * @author lizhi_duan
 * @date 2018/9/3 16:38
 */
func IsMobileLogin(ctx context.Context) (mobile bool, platform string) {
	mobile = false
	platform = MASTER_PLATFROM
	//先走一遍iris自身的验证，如果过了，直接返回，没有过，走自己的验证
	if mobile = ctx.IsMobile(); !mobile {
		//获取header中的isMobile,如果为true，代表是移动端
		isMobile := ctx.GetHeader("isMobile")
		b, e := strconv.ParseBool(isMobile)

		if e != nil {
			return
		} else {
			mobile = b
		}
		//判断平台类型
		platform = ctx.GetHeader("platform")
	} else {
		//判断平台类型
		platform = ctx.GetHeader("platform")
	}
	return
}
