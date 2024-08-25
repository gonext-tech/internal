// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package partials

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func Navbar(username string, fromProtected bool) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!-- Navbar --><nav class=\"flex items-center w-full h-24 select-none\" x-data=\"{ showMenu: false }\"><div class=\"relative flex flex-wrap items-start justify-between w-full mx-auto font-medium md:items-center md:h-24 md:justify-between\"><a href=\"/\" class=\"flex items-center w-1/4 py-4 pl-6 pr-4 space-x-2 font-extrabold text-white md:py-0\"><span class=\"flex items-center justify-center flex-shrink-0 w-8 h-8 text-gray-900 rounded-full bg-gradient-to-br from-white via-gray-200 to-white\"><img src=\"/public/logo.svg\" alt=\"site-logo\"></span> <span>Blood Dono</span></a><div :class=\"{&#39;flex&#39;: showMenu, &#39;hidden md:flex&#39;: !showMenu }\" class=\"absolute z-50 flex-col items-center justify-center w-full h-auto px-2 text-center text-gray-400 -translate-x-1/2 border-0 border-gray-700 rounded-full md:border md:w-auto md:h-10 left-1/2 md:flex-row md:items-center\"><a href=\"/\" class=\"relative inline-block w-full h-full px-4 py-5 mx-2 font-medium leading-tight text-center text-white md:py-2 group md:w-auto md:px-2 lg:mx-3 md:text-center\"><span>Home</span> <span class=\"absolute bottom-0 left-0 w-full h-px duration-300 ease-out translate-y-px bg-gradient-to-r md:from-gray-700 md:via-gray-400 md:to-gray-700 from-gray-900 via-gray-600 to-gray-900\"></span></a> <a href=\"/donor\" class=\"relative inline-block w-full h-full px-4 py-5 mx-2 font-medium leading-tight text-center duration-300 ease-out md:py-2 group hover:text-white md:w-auto md:px-2 lg:mx-3 md:text-center\"><span>Donors</span> <span class=\"absolute bottom-0 w-0 h-px duration-300 ease-out translate-y-px group-hover:left-0 left-1/2 group-hover:w-full bg-gradient-to-r md:from-gray-700 md:via-gray-400 md:to-gray-700 from-gray-900 via-gray-600 to-gray-900\"></span></a> <a href=\"/city\" class=\"relative inline-block w-full h-full px-4 py-5 mx-2 font-medium leading-tight text-center duration-300 ease-out md:py-2 group hover:text-white md:w-auto md:px-2 lg:mx-3 md:text-center\"><span>City</span> <span class=\"absolute bottom-0 w-0 h-px duration-300 ease-out translate-y-px group-hover:left-0 left-1/2 group-hover:w-full bg-gradient-to-r md:from-gray-700 md:via-gray-400 md:to-gray-700 from-gray-900 via-gray-600 to-gray-900\"></span></a> <a href=\"/blood-type\" class=\"relative inline-block w-full h-full px-4 py-5 mx-2 font-medium leading-tight text-center duration-300 ease-out md:py-2 group hover:text-white md:w-auto md:px-2 lg:mx-3 md:text-center\"><span>Blood Type</span> <span class=\"absolute bottom-0 w-0 h-px duration-300 ease-out translate-y-px group-hover:left-0 left-1/2 group-hover:w-full bg-gradient-to-r md:from-gray-700 md:via-gray-400 md:to-gray-700 from-gray-900 via-gray-600 to-gray-900\"></span></a></div><div class=\"fixed top-0 left-0 z-40 items-center hidden w-full h-full p-3 text-sm bg-gray-900 bg-opacity-50 md:w-auto md:bg-transparent md:p-0 md:relative md:flex\" :class=\"{&#39;flex&#39;: showMenu, &#39;hidden&#39;: !showMenu }\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !fromProtected {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"flex-col items-center w-full h-full p-3 overflow-hidden bg-black bg-opacity-50 rounded-lg select-none md:p-0 backdrop-blur-lg md:h-auto md:bg-transparent md:rounded-none md:relative md:flex md:flex-row md:overflow-auto\"><div class=\"flex flex-col items-center justify-end w-full h-full pt-2 md:w-full md:flex-row md:py-0\"><a href=\"/login\" class=\"w-full py-5 mr-0 text-center text-gray-200 md:py-3 md:w-auto hover:text-white md:pl-0 md:mr-3 lg:mr-5\">Sign In</a> <a href=\"/register\" class=\"inline-flex items-center justify-center w-full px-4 py-3 md:py-1.5 font-medium leading-6 text-center whitespace-no-wrap transition duration-150 ease-in-out border border-transparent md:mr-1 text-gray-600 md:w-auto bg-white rounded-lg md:rounded-full hover:bg-white focus:outline-none focus:border-gray-700 focus:shadow-outline-gray active:bg-gray-700\">Sign Up</a></div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div @click=\"showMenu = !showMenu\" class=\"absolute right-0 z-50 flex flex-col items-end translate-y-1.5 w-10 h-10 p-2 mr-4 rounded-full cursor-pointer md:hidden hover:bg-gray-200/10 hover:bg-opacity-10\" :class=\"{ &#39;text-gray-400&#39;: showMenu, &#39;text-gray-100&#39;: !showMenu }\"><svg class=\"w-6 h-6\" x-show=\"!showMenu\" fill=\"none\" stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" viewBox=\"0 0 24 24\" stroke=\"currentColor\" x-cloak><path d=\"M4 6h16M4 12h16M4 18h16\"></path></svg> <svg class=\"w-6 h-6\" x-show=\"showMenu\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\" xmlns=\"http://www.w3.org/2000/svg\" x-cloak><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M6 18L18 6M6 6l12 12\"></path></svg></div></div></nav><script>\n  function removeMask() {\n    //   document.getElementById(\"notification-container\").querySelector(\".mask\").classList.add(\"hidden\");\n  }\n  function fetchNotifications() {\n    const targetElement = document.getElementById(\"notification-results\");\n    if (targetElement) {\n      targetElement.setAttribute(\"hx-get\", \"/notification/navbar\");\n      targetElement.setAttribute(\"hx-target\", \"#notification-results\");\n      targetElement.setAttribute(\"hx-swap\", \"outer\"); // Optional: Swap entire content\n      const url = targetElement.getAttribute(\"hx-get\");\n\n      // Make the HTTP request (replace with your preferred method)\n      fetch(url)\n        .then(response => response.text())\n        .then(data => {\n          targetElement.innerHTML = data; // Update content (adjust based on response format)\n        })\n        .catch(error => {\n          console.error(\"Error fetching notifications:\", error);\n          // Implement error handling (optional)\n        });\n    }\n    //   document.getElementById(\"notification-container\").querySelector(\".mask\").classList.remove(\"hidden\");\n  }\n  var loc = window.location;\n  var uri = 'ws:';\n  if (loc.protocol === 'https:') {\n    uri = 'wss:';\n  }\n  uri += '//' + loc.host;\n  uri += '/ws/ticket'\n\n  ws = new WebSocket(uri)\n\n  ws.onopen = function () {\n    console.log('Connected')\n  }\n\n  ws.onmessage = function (evt) {\n    console.log('message', evt.data)\n    if (evt.data === \"refetch\") {\n      fetchNotifications()\n    }\n  }\n\n</script><script>\n\t//  let currentTime = new Date();\n\t// const notificationItems = document.querySelectorAll('#notifications li[data-value]');\n\t//  console.log(\"notifications\", notificationItems);\n\t//  notificationItems.forEach((item) => {\n\t//   // Get the value of the data-value attribute\n\t//   const timeValue = item.getAttribute(\"data-value\");\n\t//   console.log(\"timeValue\", timeValue);\n\t//   // Parse the time value into a Date object\n\t//   const notificationTime = new Date(timeValue);\n\t//   // Calculate the time difference in milliseconds\n\t//   const timeDiff = currentTime - notificationTime;\n\t//   // Convert milliseconds to minutes\n\t//   const minutesDiff = Math.floor(timeDiff / (1000 * 60));\n\t//   // Display the time difference as \"x minutes ago\"\n\t//   if (minutesDiff < 60) {\n\t//    item.textContent = `${minutesDiff}m ago`;\n\t//   } else {\n\t//    // Convert minutes to hours\n\t//    const hoursDiff = Math.floor(minutesDiff / 60);\n\t//    item.textContent = `${hoursDiff}h ago`;\n\t//   }\n\t//  });\n\t// </script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
