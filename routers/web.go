package routers

import (
	"github.com/astaxie/beego"
	"programming-learning-platform/controllers"
)

func webRouter() {
	// 基本
	beego.Router("/crawl", &controllers.BaseController{}, "post:Crawl")
	beego.Router("/sitemap.html", &controllers.BaseController{}, "get:Sitemap") // 站点地图

	// 静态文件
	beego.Router("/app", &controllers.StaticController{}, "get:APP")
	beego.Router("/projects/*", &controllers.StaticController{}, "get:ProjectsFile")
	beego.Router("/uploads/*", &controllers.StaticController{}, "get:Uploads")
	beego.Router("/*", &controllers.StaticController{}, "get:StaticFile")

	// 登陆注册页
	beego.Router("/login", &controllers.AccountController{}, "*:Login")                // 用户普通登陆
	beego.Router("/login/:oauth", &controllers.AccountController{}, "*:Oauth")         // 用户第三方登陆鉴权
	beego.Router("/logout", &controllers.AccountController{}, "*:Logout")              // 用户退出登陆
	beego.Router("/bind", &controllers.AccountController{}, "post:Bind")               // 用户注册
	beego.Router("/find_password", &controllers.AccountController{}, "*:FindPassword") // 找回密码
	beego.Router("/valid_email", &controllers.AccountController{}, "post:ValidEmail")  // 邮箱校验

	// 首页
	beego.Router("/", &controllers.CateController{}, "get:Index")
	beego.Router("/cate", &controllers.CateController{}, "get:List")
	// 发现页
	beego.Router("/explore", &controllers.HomeController{}, "*:Index")
	// 榜单页
	beego.Router("/rank", &controllers.RankController{}, "get:Index")
	// 标签页
	beego.Router("/label", &controllers.LabelController{}, "get:List")
	beego.Router("/label/:key", &controllers.LabelController{}, "get:Index")
	// 友情链接页
	beego.Router("/friendlink", &controllers.FriendLinkController{}, "*:List")
	// 我的收藏
	beego.Router("/star", &controllers.StarController{}, "*:List")
	// 书籍搜索
	beego.Router("/search", &controllers.SearchController{}, "get:Search")        // 搜索页
	beego.Router("/search/result", &controllers.SearchController{}, "get:Result") // 搜索结果

	// 个人主页
	beego.Router("/user/:username", &controllers.UserController{}, "get:Index")                 // 个人首页
	beego.Router("/user/:username/collection", &controllers.UserController{}, "get:Collection") // 个人收藏
	beego.Router("/user/:username/follow", &controllers.UserController{}, "get:Follow")         // 个人关注
	beego.Router("/user/:username/fans", &controllers.UserController{}, "get:Fans")             // 个人粉丝
	beego.Router("/user/sign", &controllers.BaseController{}, "get:SignToday")                  // 每日签到
	beego.Router("/follow/:uid", &controllers.BaseController{}, "get:SetFollow")                // 关注或取消关注

	// 个人设置
	beego.Router("/setting", &controllers.SettingController{}, "*:Index")             // 个人设置
	beego.Router("/setting/password", &controllers.SettingController{}, "*:Password") // 修改密码
	beego.Router("/setting/upload", &controllers.SettingController{}, "*:Upload")     // 上传图片
	beego.Router("/setting/qrcode", &controllers.SettingController{}, "*:Qrcode")

	// 我的书籍
	beego.Router("/book", &controllers.BookController{}, "*:Index")                              // 书籍列表
	beego.Router("/book/create", &controllers.BookController{}, "post:Create")                   // 添加书籍
	beego.Router("/book/star/:id", &controllers.BookController{}, "*:Star")                      // 收藏书籍
	beego.Router("/book/score/:id", &controllers.BookController{}, "*:Score")                    // 书籍评分
	beego.Router("/book/comment/:id", &controllers.BookController{}, "post:Comment")             // 书籍评论
	beego.Router("/book/uploadProject", &controllers.BookController{}, "post:UploadProject")     // zip、EPUB上传导入书籍
	beego.Router("/book/downloadProject", &controllers.BookController{}, "post:DownloadProject") // zip拉取导入书籍
	beego.Router("/book/git-pull", &controllers.BookController{}, "post:GitPull")                // git clone导入书籍
	beego.Router("/book/:key/dashboard", &controllers.BookController{}, "*:Dashboard")           // 书籍概要
	beego.Router("/book/:key/setting", &controllers.BookController{}, "*:Setting")               // 书籍设置
	beego.Router("/book/:key/users", &controllers.BookController{}, "*:Users")                   // 书籍成员
	beego.Router("/book/:key/release", &controllers.BookController{}, "post:Release")            // 书籍发布
	beego.Router("/book/:key/generate", &controllers.BookController{}, "get,post:Generate")      // 生成下载文档
	beego.Router("/book/:key/sort", &controllers.BookController{}, "post:SaveSort")              // 文档排序
	beego.Router("/book/:key/replace", &controllers.BookController{}, "get,post:Replace")        //
	beego.Router("/book/setting/save", &controllers.BookController{}, "post:SaveBook")           // 保存书籍修改
	beego.Router("/book/setting/private", &controllers.BookController{}, "post:PrivatelyOwned")  // 设置书籍私有状态
	beego.Router("/book/setting/transfer", &controllers.BookController{}, "post:Transfer")       // 书籍转让
	beego.Router("/book/setting/upload", &controllers.BookController{}, "post:UploadCover")      // 上传书籍封面
	beego.Router("/book/setting/token", &controllers.BookController{}, "post:CreateToken")       // 生成私有书籍阅读令牌
	beego.Router("/book/setting/delete", &controllers.BookController{}, "post:Delete")           // 删除书籍
	beego.Router("/book/users/create", &controllers.BookMemberController{}, "post:AddMember")    // 添加书籍成员
	beego.Router("/book/users/change", &controllers.BookMemberController{}, "post:ChangeRole")   // 修改书籍成员角色
	beego.Router("/book/users/delete", &controllers.BookMemberController{}, "post:RemoveMember") // 删除书籍成员

	// 管理后台(管理员用户)
	beego.Router("/manager", &controllers.ManagerController{}, "*:Index")                                              // 管理后台首页
	beego.Router("/manager/users", &controllers.ManagerController{}, "*:Users")                                        // 用户列表
	beego.Router("/manager/users/edit/:id", &controllers.ManagerController{}, "*:EditMember")                          // 修改用户信息
	beego.Router("/manager/member/create", &controllers.ManagerController{}, "post:CreateMember")                      // 添加用户
	beego.Router("/manager/member/delete", &controllers.ManagerController{}, "post:DeleteMember")                      // 删除用户
	beego.Router("/manager/member/update-member-status", &controllers.ManagerController{}, "post:UpdateMemberStatus")  // 更新用户状态
	beego.Router("/manager/member/update-member-no-rank", &controllers.ManagerController{}, "post:UpdateMemberNoRank") // 更新用户是否排榜
	beego.Router("/manager/member/change-member-role", &controllers.ManagerController{}, "post:ChangeMemberRole")      // 更新用户角色
	beego.Router("/manager/books", &controllers.ManagerController{}, "*:Books")                                        // 书籍列表
	beego.Router("/manager/books/edit/:key", &controllers.ManagerController{}, "*:EditBook")                           // 书籍设置
	beego.Router("/manager/books/delete", &controllers.ManagerController{}, "*:DeleteBook")                            // 删除书籍
	beego.Router("/manager/books/transfer", &controllers.ManagerController{}, "post:Transfer")                         // 转让书籍
	beego.Router("/manager/books/sort", &controllers.ManagerController{}, "get:UpdateBookSort")                        // 更新书籍排序
	beego.Router("/manager/books/open", &controllers.ManagerController{}, "post:PrivatelyOwned")                       // 设置书籍私有状态
	beego.Router("/manager/books/token", &controllers.ManagerController{}, "post:CreateToken")                         // 创建令牌
	beego.Router("/manager/comments", &controllers.ManagerController{}, "*:Comments")                                  // 评论列表
	beego.Router("/manager/comments/delete", &controllers.ManagerController{}, "*:DeleteComment")                      // 删除评论
	beego.Router("/manager/comments/clear", &controllers.ManagerController{}, "*:ClearComments")                       // 清除当前用户评论
	beego.Router("/manager/comments/set", &controllers.ManagerController{}, "*:SetCommentStatus")                      // 设置评论状态
	beego.Router("/manager/setting", &controllers.ManagerController{}, "*:Setting")                                    // 配置管理
	beego.Router("/manager/attach/list", &controllers.ManagerController{}, "*:AttachList")                             // 附件列表
	beego.Router("/manager/attach/detailed/:id", &controllers.ManagerController{}, "*:AttachDetailed")                 // 附件详情
	beego.Router("/manager/attach/delete", &controllers.ManagerController{}, "*:AttachDelete")                         // 删除附件
	beego.Router("/manager/tags", &controllers.ManagerController{}, "get:Tags")                                        // 标签列表
	beego.Router("/manager/add-tags", &controllers.ManagerController{}, "post:AddTags")                                // 添加标签
	beego.Router("/manager/del-tags", &controllers.ManagerController{}, "get:DelTags")                                 // 删除标签
	beego.Router("/manager/seo", &controllers.ManagerController{}, "post,get:Seo")                                     // seo管理
	beego.Router("/manager/ads", &controllers.ManagerController{}, "post,get:Ads")                                     // 广告管理
	beego.Router("/manager/update-ads", &controllers.ManagerController{}, "post,get:UpdateAds")                        // 修改广告信息
	beego.Router("/manager/del-ads", &controllers.ManagerController{}, "get:DelAds")                                   // 删除广告
	beego.Router("/manager/category", &controllers.ManagerController{}, "post,get:Category")                           // 分类列表
	beego.Router("/manager/update-cate", &controllers.ManagerController{}, "get:UpdateCate")                           // 更新分类
	beego.Router("/manager/del-cate", &controllers.ManagerController{}, "get:DelCate")                                 // 删除分类
	beego.Router("/manager/icon-cate", &controllers.ManagerController{}, "post:UpdateCateIcon")                        // 更新分类的图标
	beego.Router("/manager/icon-friendlink", &controllers.ManagerController{}, "post:UpdateFriendLinkIcon")            // 更新友链的图标
	beego.Router("/manager/sitemap", &controllers.ManagerController{}, "get:Sitemap")                                  // 更新站点地图
	beego.Router("/manager/friendlink", &controllers.ManagerController{}, "get:FriendLink")                            // 友链管理
	beego.Router("/manager/add_friendlink", &controllers.ManagerController{}, "post:AddFriendLink")                    // 添加友链
	beego.Router("/manager/update_friendlink", &controllers.ManagerController{}, "get:UpdateFriendLink")               // 更新友链
	beego.Router("/manager/del_friendlink", &controllers.ManagerController{}, "get:DelFriendLink")                     // 删除友链
	beego.Router("/manager/rebuild-index", &controllers.ManagerController{}, "get:RebuildAllIndex")                    // 重建全量索引
	beego.Router("/manager/banners", &controllers.ManagerController{}, "get:Banners")                                  // 横幅管理
	beego.Router("/manager/banners/upload", &controllers.ManagerController{}, "post:UploadBanner")                     // 上传横幅
	beego.Router("/manager/banners/delete", &controllers.ManagerController{}, "get:DeleteBanner")                      // 删除横幅
	beego.Router("/manager/banners/update", &controllers.ManagerController{}, "get:UpdateBanner")                      // 更新横幅

	// 书签
	beego.Router("/bookmark/:id", &controllers.BookmarkController{}, "get:Bookmark")
	beego.Router("/bookmark/list/:book_id", &controllers.BookmarkController{}, "get:List")

	// 阅读记录
	beego.Router("/record/:book_id", &controllers.RecordController{}, "get:List")
	beego.Router("/record/:book_id/clear", &controllers.RecordController{}, "get:Clear")
	beego.Router("/record/delete/:doc_id", &controllers.RecordController{}, "get:Delete")

	// 文档
	beego.Router("/api/attach/remove/", &controllers.DocumentController{}, "post:RemoveAttachment")
	beego.Router("/api/:key/edit/?:id", &controllers.DocumentController{}, "*:Edit")
	beego.Router("/api/upload", &controllers.DocumentController{}, "post:Upload")
	beego.Router("/api/:key/create", &controllers.DocumentController{}, "post:Create")
	beego.Router("/api/create_multi", &controllers.DocumentController{}, "post:CreateMulti")
	beego.Router("/api/:key/delete", &controllers.DocumentController{}, "post:Delete")
	beego.Router("/api/:key/content/?:id", &controllers.DocumentController{}, "*:Content")
	beego.Router("/api/:key/compare/:id", &controllers.DocumentController{}, "*:Compare")
	beego.Router("/history/get", &controllers.DocumentController{}, "get:History")
	beego.Router("/history/delete", &controllers.DocumentController{}, "*:DeleteHistory")
	beego.Router("/history/restore", &controllers.DocumentController{}, "*:RestoreHistory")
	beego.Router("/books/:key", &controllers.DocumentController{}, "*:Index") // 文档简介
	beego.Router("/read/:key/:id", &controllers.DocumentController{}, "*:Read")
	beego.Router("/read/:key/search", &controllers.DocumentController{}, "post:Search")
	beego.Router("/export/:key", &controllers.DocumentController{}, "*:Export")
	beego.Router("/qrcode/:key.png", &controllers.DocumentController{}, "get:QrCode")
	beego.Router("/attach_files/:key/:attach_id", &controllers.DocumentController{}, "get:DownloadAttachment")

	beego.Router("/local-render", &controllers.LocalhostController{}, "get,post:RenderMarkdown")
	beego.Router("/local-render-cover", &controllers.LocalhostController{}, "get:RenderCover")
}
