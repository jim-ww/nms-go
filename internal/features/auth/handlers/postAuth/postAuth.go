package postauth

//
// type AuthHandler struct {
// 	authService *authService.AuthService
// 	jwtService  *jwtService.JWTService
// 	logger      *slog.Logger
// 	tmpl        *template.Template
// }
//
// func New(log *slog.Logger, userService *authService.AuthService, jwtService *jwtService.JWTService, tmplHandler *handlers.TmplHandler) *AuthHandler {
// 	templ := template.Must(template.ParseFiles("web/templates/auth.html"))
// 	return &AuthHandler{
// 		logger:      log,
// 		authService: userService,
// 		jwtService:  jwtService,
// 		tmpl:        templ,
// 	}
// }
//
// func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
// 	// parse form & convert to dtos
// 	if err := r.ParseForm(); err != nil {
// 		ah.logger.Error("unable to parse form", sl.Err(err))
// 		handlers.RenderError(w, r, "unable to parse form", http.StatusBadRequest)
// 		return
// 	}
// 	dto := dtos.NewLoginDTO(r.FormValue("username"), r.FormValue("password"))
// 	_ = dto // TODO
// 	// attempt to login
// 	// token, validationErrors, err := ah.authService.LoginUser(dto)
//
// 	// if has errors, return validation errors to form
//
// 	// if no errors, set token cookie and redirect to home page
// 	// http.SetCookie(w, ah.jwtService.NewTokenCookie(token))
// 	// http.Redirect(w, r, "/", http.StatusSeeOther)
// }
//
// func (lh *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
// 	if err := r.ParseForm(); err != nil {
// 		lh.logger.Error("Unable to parse form", sl.Err(err))
// 		handlers.RenderError(w, r, "Unable to parse form", http.StatusBadRequest)
// 		return
// 	}
// 	// dto := dtos.NewRegisterDTO(r.FormValue("username"), r.FormValue("email"), r.FormValue("password"))
//
// 	// token, validationErrors, err := lh.authService.RegisterUser(dto)
//
// 	// if has errors, return validation errors to form
//
// 	// if no errors, set token cookie and redirect to home page
// 	// http.SetCookie(w, lh.jwtService.NewTokenCookie(token))
// 	// http.Redirect(w, r, "/", http.StatusSeeOther)
// }
