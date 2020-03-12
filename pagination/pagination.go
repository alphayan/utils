package pagination

import "github.com/astaxie/beego/orm"

type Pagination struct {
	CurrentPage int64 `orm:"-"` //当前页数
	PageSize    int64 `orm:"-"` //每页大小
	TotalsCount int64 `orm:"-"` //总个数
	TotalsPage  int64 `orm:"-"` //总页数
	Results     interface{}     //查询数据
}

//初始化分页对象，前台如果没有指定每页展示数量，默认每页100条
func NewPagination() (page Pagination) {
	page = Pagination{
		PageSize: 100,
		Results:  make([]interface{}, 0),
	}
	return
}

//创建orm的builder
func (page *Pagination) NewMysqlQueryBuild() (queryBuilder orm.QueryBuilder) {
	var err error
	queryBuilder, err = orm.NewQueryBuilder("mysql")
	if err != nil {
		panic(err)
	}
	return
}

//判断是否需要分页
func (page *Pagination) IsPagination() bool {
	if page.CurrentPage != 0 && page.PageSize != 0 {
		return true
	}
	return false
}

//分页limit索引
func (page *Pagination) Limit() int {
	return int(page.PageSize)
}

//分页开始索引
func (page *Pagination) Offset() int {
	return int((page.CurrentPage - 1) * page.PageSize)
}

//最后查询结果处理,如果不需要分页的话count传nil
func (page *Pagination) Result(count interface{}, result interface{}) {
	if count != nil {
		page.TotalsCount = count.(int64)
		if page.IsPagination() {
			if page.TotalsCount%int64(page.PageSize) == 0 {
				page.TotalsPage = page.TotalsCount / int64(page.PageSize)
			} else {
				page.TotalsPage = (page.TotalsCount / int64(page.PageSize)) + 1
			}
			if page.CurrentPage > page.TotalsPage {
				page.Results = make([]interface{}, 0)
				return
			}
		}
	}
	page.Results = result
}
