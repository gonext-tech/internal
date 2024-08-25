// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package layout

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "github.com/gonext-tech/internal/views/partials"

func Base(title, username string, fromProtected, isError bool, errMsgs, sucMsgs []string) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html lang=\"en\" data-theme=\"dark\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><meta name=\"description\" content=\"Blood donation app\"><meta name=\"google\" content=\"notranslate\"><link rel=\"shortcut icon\" href=\"/logo.svg\" type=\"image/svg\"><script src=\"https://cdn.tailwindcss.com\"></script><!-- <link rel=\"stylesheet\" href=\"/css/styles.css\"/> --><link rel=\"stylesheet\" href=\"https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.2/css/all.min.css\" integrity=\"sha512-SnH5WK+bZxgPHs44uWIX+LLJAJ9/2PkPKZ5QiAj6Ta86w+fsb2TkcmfRyVX3pBnMFcV7oQPJkl9QevSCWr3W6A==\" crossorigin=\"anonymous\" referrerpolicy=\"no-referrer\"><link rel=\"stylesheet\" href=\"/css/main.css\"><link href=\"https://cdn.jsdelivr.net/npm/daisyui@4.4.10/dist/full.min.css\" rel=\"stylesheet\" type=\"text/css\"><title>Internal | ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(title)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/layout/base.layout.templ`, Line: 24, Col: 28}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</title><script src=\"https://unpkg.com/htmx.org@1.9.9\" integrity=\"sha384-QFjmbokDn2DjBjq+fM+8LUIVrAgqcNW2s0PjAxHETgRn9l4fvX31ZxDxvwQnyMOX\" crossorigin=\"anonymous\"></script><script src=\"https://unpkg.com/htmx.org/dist/ext/ws.js\"></script><script src=\"https://unpkg.com/hyperscript.org@0.9.12\"></script></head><body id=\"main-content\" class=\"sample-transition bg-base-100\" hx-boost=\"true\"><div class=\"flex\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !isError && len(username) > 0 {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<aside class=\"sticky top-0 h-screen w-64\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = partials.Sidebar(username, fromProtected).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</aside>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"flex-1\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 = []any{templ.KV("pt-10 px-10 min-h-screen w-full", !isError)}
		templ_7745c5c3_Err = templ.RenderCSSItems(ctx, templ_7745c5c3_Buffer, templ_7745c5c3_Var3...)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<main class=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ.CSSClasses(templ_7745c5c3_Var3).String()))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ_7745c5c3_Var1.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div className=\"relative\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = partials.FlashMessages(errMsgs, sucMsgs).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></main><footer class=\"w-full\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = partials.Footer().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</footer></div></div></body><script>\n\n  // --> ANITMATION START <--\n\n  document.addEventListener('htmx:beforeSwap', function (event) {\n    // Check if the swap is coming from an element with a specific data attribute\n\n    if (event.detail.requestConfig.headers[\"HX-Trigger\"] === 'back_button') {\n      // Apply custom logic for the back button swap\n      console.log('Navigating back with custom logic.');\n\n      // Example: Change the transition class\n      const mainContent = document.getElementById('main-content');\n      if (mainContent) {\n        mainContent.classList.remove('sample-transition');\n        mainContent.classList.add('sample-transition-back');\n      }\n    } else {\n      // Handle other cases (optional)\n      console.log('Navigating forward or from other elements.');\n\n      // Reset to the default transition\n      const mainContent = document.getElementById('main-content');\n      if (mainContent) {\n        mainContent.classList.remove('sample-transition-back');\n        mainContent.classList.add('sample-transition');\n      }\n    }\n  });\n\n  window.addEventListener('popstate', function () {\n    const container = document.getElementById('main-content');\n    if (container) {\n      container.classList.add('sample-transition');\n      setTimeout(function () {\n        container.classList.remove('sample-transition');\n      }, 600);\n    }\n  });\n\n  // --> ANITMATION END <--\n\n  function confirmDelete(id, path) {\n    if (confirm(\"Are you sure you want to delete this item?\")) {\n      const deleteUrl = `/${path}/${id}`;\n\n      // Create a temporary element to use HTMX attributes\n      const tempElement = document.createElement('div');\n      tempElement.setAttribute('hx-delete', deleteUrl);\n      tempElement.setAttribute('hx-swap', 'transition:true');\n      tempElement.setAttribute('hx-target', 'body');\n\n      htmx.ajax('DELETE', deleteUrl, {\n        target: 'body',\n        swap: 'transition:true'\n      }).then(() => {\n        alert('Item deleted successfully');\n      }).catch(err => {\n        console.error(err);\n        alert('Error deleting item');\n      });\n    }\n  }\n</script></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
