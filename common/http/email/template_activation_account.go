package email

import (
	"fmt"

	"github.com/spf13/viper"
)

func GenerateTemplateActivationAccount(linkActivation string, userName string, userEmail string, config *viper.Viper) (html string) {
	dataConfig := config.Sub("general")

	headerText := fmt.Sprintf(`Hello %s<%s>, <br>Thank you for Registration!`, userName, userEmail)
	bodyText := fmt.Sprintf(`You are just one step away from completing your registration, activate your account by clicking the button below to start your jurney in %s and gain access to Clock In and Clock Out!`, dataConfig.GetString("company_name"))

	html = fmt.Sprintf(`
	<!-- START HEAD -->
   <head>
   <!-- CHARSET -->
   <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
   <!-- MOBILE FIRST -->
   <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0">
   <!-- GOOGLE FONTS -->
   <link href="https://fonts.googleapis.com/css?family=Ubuntu+Mono" rel="stylesheet">
   <link href="https://fonts.googleapis.com/css?family=Ubuntu" rel="stylesheet">
   <!-- RESPONSIVE CSS -->
   <style type="text/css">
      @media only screen and (max-width: 550px){
      .responsive_at_550{
      width: 90%% !important;
      max-width: 90%% !important;
      }
      }
   </style>
   </head>
   <!-- END HEAD -->
   <!-- START BODY -->
   <body leftmargin="0" topmargin="0" marginwidth="0" marginheight="0">
      <!-- START EMAIL CONTENT -->
      <table width="100%%" border="0" cellpadding="0" cellspacing="0" align="center">
         <tbody>
            <tr>
               <td align="center" bgcolor="#1976D2">
                  <table width="100%%" border="0" cellpadding="0" cellspacing="0" align="center">
                     <tbody>
                        <tr>
                           <td width="100%%" align="center">
                              <!-- START SPACING -->
                              <table width="100%%" border="0" cellpadding="0" cellspacing="0" align="center">
                                 <tbody>
                                    <tr>
                                       <td height="40">&nbsp;</td>
                                    </tr>
                                 </tbody>
                              </table>
                              <!-- END SPACING -->
                              <!-- START LOGO -->
                              <table width="200" border="0" cellpadding="0" cellspacing="0" align="center">
                                 <tbody>
                                    <tr>
                                       <td width="100%%" align="center">
                                          <img width="200" src="%s" alt="Wokdev" border="0" style="text-align: center;"/>
                                       </td>
                                    </tr>
                                 </tbody>
                              </table>
                              <!-- END LOGO -->
                              <!-- START SPACING -->
                              <table width="100%%" border="0" cellpadding="0" cellspacing="0" align="center">
                                 <tbody>
                                    <tr>
                                       <td height="40">&nbsp;</td>
                                    </tr>
                                 </tbody>
                              </table>
                              <!-- END SPACING -->
                              <!-- START CONTENT -->
                              <table width="500" border="0" cellpadding="0" cellspacing="0" align="center" style="padding-left:20px; padding-right:20px;" class="responsive_at_550">
                                 <tbody>
                                    <tr>
                                       <td align="center" bgcolor="#ffffff">
                                          <!-- START BORDER COLOR -->
                                          <table width="100%%" border="0" cellpadding="0" cellspacing="0" align="center">
                                             <tbody>
                                                <tr>
                                                   <td width="100%%" height="7" align="center" border="0" bgcolor="#03a9f4"></td>
                                                </tr>
                                             </tbody>
                                          </table>
                                          <!-- END BORDER COLOR -->
                                          <!-- START SPACING -->
                                          <table width="100%%" border="0" cellpadding="0" cellspacing="0" align="center">
                                             <tbody>
                                                <tr>
                                                   <td height="30">&nbsp;</td>
                                                </tr>
                                             </tbody>
                                          </table>
                                          <!-- END SPACING -->
                                          <!-- START HEADING -->
                                          <table width="90%%" border="0" cellpadding="0" cellspacing="0" align="center">
                                             <tbody>
                                                <tr>
                                                   <td width="100%%" align="center">
                                                      <h1 style="font-family:'Ubuntu Mono', monospace; font-size:20px; color:#202020; font-weight:bold; padding-left:20px; padding-right:20px;">
                                                      %s
                                                      </h1>
                                                   </td>
                                                </tr>
                                             </tbody>
                                          </table>
                                          <!-- END HEADING -->
                                          <!-- START PARAGRAPH -->
                                          <table width="90%%" border="0" cellpadding="0" cellspacing="0" align="center">
                                             <tbody>
                                                <tr>
                                                   <td width="100%%" align="center">
                                                      <p style="font-family:'Ubuntu', sans-serif; font-size:14px; color:#202020; padding-left:20px; padding-right:20px; text-align:justify;">
                                                      %s
                                                      </p>
                                                   </td>
                                                </tr>
                                             </tbody>
                                          </table>
                                          <!-- END PARAGRAPH -->
                                          <!-- START SPACING -->
                                          <table width="100%%" border="0" cellpadding="0" cellspacing="0" align="center">
                                             <tbody>
                                                <tr>
                                                   <td height="30">&nbsp;</td>
                                                </tr>
                                             </tbody>
                                          </table>
                                          <!-- END SPACING -->
                                          %s
                                          <!-- START BUTTON -->
                                          <table width="200" border="0" cellpadding="0" cellspacing="0" align="center">
                                             <tbody>
                                                <tr>
                                                   <td align="center" bgcolor="#1976D2">
                                                      <a style="font-family:'Ubuntu Mono', monospace; display:block; color:#ffffff; font-size:14px; font-weight:bold; text-decoration:none; padding-left:20px; padding-right:20px; padding-top:20px; padding-bottom:20px;" href="%s">Verify E-mail Address</a>
                                                   </td>
                                                </tr>
                                             </tbody>
                                          </table>
                                          <!-- END BUTTON -->
                                          <!-- START SPACING -->
                                          <table width="100%%" border="0" cellpadding="0" cellspacing="0" align="center">
                                             <tbody>
                                                <tr>
                                                   <td height="30">&nbsp;</td>
                                                </tr>
                                             </tbody>
                                          </table>
                                          <!-- END SPACING -->
                                       </td>
                                    </tr>
                                 </tbody>
                              </table>
                              <!-- END CONTENT -->
                              <!-- START SPACING -->
                              <table width="100%%" border="0" cellpadding="0" cellspacing="0" align="center">
                                 <tbody>
                                    <tr>
                                       <td height="40">&nbsp;</td>
                                    </tr>
                                 </tbody>
                              </table>
                              <!-- END SPACING -->
                              <!-- START SOCIAL MEDIA ICONS -->
                              <table width="100%%" border="0" cellpadding="0" cellspacing="0" align="center">
                                 <tbody>
                                    <tr>
                                       <td width="100%%" align="center">
                                          <a href="%s"><img width="25" height="25" src="https://codewiz.co/img/email_templates/facebook_icon.png" alt="Facebook" border="0" style="text-align: center;"/></a>
                                          <a href="%s"><img width="25" height="25" src="https://codewiz.co/img/email_templates/twitter_icon.png" alt="Twitter" border="0" style="text-align: center;"/></a>
                                          <a href="%s"><img width="25" height="25" src="https://codewiz.co/img/email_templates/linkedin_icon.png" alt="LinkedIn" border="0" style="text-align: center;"/></a>
                                          <a href="%s"><img width="25" height="25" src="https://codewiz.co/img/email_templates/instagram_icon.png" alt="Instagram" border="0" style="text-align: center;"/></a>
                                          <a href="%s"><img width="25" height="25" src="https://codewiz.co/img/email_templates/youtube_icon.png" alt="Youtube" border="0" style="text-align: center;"/></a>
                                          <a href="%s"><img width="25" height="25" src="https://codewiz.co/img/email_templates/google_plus_icon.png" alt="Google Plus" border="0" style="text-align: center;"/></a>
                                          <a href="%s"><img width="25" height="25" src="https://codewiz.co/img/email_templates/github_icon.png" alt="Github" border="0" style="text-align: center;"/></a>
                                       </td>
                                    </tr>
                                 </tbody>
                              </table>
                              <!-- END SOCIAL MEDIA ICONS -->
                              <!-- START FOOTER -->
                              <table width="100%%" border="0" cellpadding="0" cellspacing="0" align="center">
                                 <tbody>
                                    <tr>
                                       <td width="100%%" align="center" style="padding-left:15px; padding-right:15px;">
                                          <p style="font-family:'Ubuntu Mono', monospace; color:#ffffff; font-size:12px;">%s &copy; 2017, All Rights Reserved</p>
                                       </td>
                                    </tr>
                                    <tr>
                                       <td width="100%%" align="center" style="padding-left:15px; padding-right:15px;">
                                          <a href="%s" style="text-decoration:underline; font-family:'Ubuntu Mono', monospace; color:#ffffff; font-size:12px;">Terms of Use</a>
                                          <span style="font-family:'Ubuntu Mono', monospace; color:#ffffff;">|</span>
                                          <a href="%s" style="text-decoration:underline; font-family:'Ubuntu Mono', monospace; color:#ffffff; font-size:12px;">Privacy Policy</a>
                                       </td>
                                    </tr>
                                 </tbody>
                              </table>
                              <!-- END FOOTER -->
                              <!-- START SPACING -->
                              <table width="100%%" border="0" cellpadding="0" cellspacing="0" align="center">
                                 <tbody>
                                    <tr>
                                       <td height="40">&nbsp;</td>
                                    </tr>
                                 </tbody>
                              </table>
                              <!-- END SPACING -->
                           </td>
                        </tr>
                     </tbody>
                  </table>
               </td>
            </tr>
         </tbody>
      </table>
      <!-- END EMAIL CONTENT -->
   </body>
   <!-- END BODY -->`,
		dataConfig.GetString("app_logo"),
		headerText,
		bodyText,
		linkActivation,
		linkActivation,
		dataConfig.GetString("facebook"),
		dataConfig.GetString("twitter"),
		dataConfig.GetString("linkedin"),
		dataConfig.GetString("instagram"),
		dataConfig.GetString("youtube"),
		dataConfig.GetString("email"),
		dataConfig.GetString("github"),
		dataConfig.GetString("app_name"),
		dataConfig.GetString("term_of_use"),
		dataConfig.GetString("privacy_policy"),
	)

	return
}
