// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package project_views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "github.com/gonext-tech/internal/views/components"

func Create() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<section class=\" bg-base-200 rounded-md relative\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.BackButton("project").Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"py-8 px-4 mx-auto max-w-2xl lg:py-16\"><h2 class=\"mb-4 text-xl font-bold text-gray-900 dark:text-white\">Add new project</h2><form action=\"\" method=\"post\" hx-swap=\"transition:true\" hx-encoding=\"multipart/form-data\"><!--   Image upload START -->")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.UploadImage("", "shop").Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"grid gap-4 sm:grid-cols-2 sm:gap-6\"><div class=\"sm:col-span-2\"><label for=\"name\" class=\"block mb-2 text-sm font-medium text-gray-900 dark:text-white\">Name</label> <input type=\"text\" name=\"name\" id=\"name\" class=\"bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-primary-500 dark:focus:border-primary-500\" placeholder=\"Cutest shop\" required=\"\"></div><div class=\"\"><label for=\"db_name\" class=\"block mb-2 text-sm font-medium text-gray-900 dark:text-white\">DB Name</label> <input type=\"text\" name=\"db_name\" id=\"db_name\" class=\"bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-primary-500 dark:focus:border-primary-500\" placeholder=\"QwikDB\" required=\"\"></div><div class=\"\"><label for=\"domain_url\" class=\"block mb-2 text-sm font-medium text-gray-900 dark:text-white\">Domain</label> <input type=\"text\" name=\"domain_url\" id=\"domain_url\" class=\"bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-primary-500 dark:focus:border-primary-500\" placeholder=\"https://qwik.gonext.tech\" required=\"\"></div><div class=\"\"><label for=\"repo_name\" class=\"block mb-2 text-sm font-medium text-gray-900 dark:text-white\">Repo Name</label> <input type=\"text\" name=\"repo_name\" id=\"repo_name\" class=\"bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-primary-500 dark:focus:border-primary-500\" placeholder=\"Qwik\" required=\"\"></div><div class=\"\"><label for=\"ssl_expired_at\" class=\"block mb-2 text-sm font-medium text-gray-900 dark:text-white\">SSL Expiration</label> <input type=\"date\" name=\"ssl_expired_at\" id=\"ssl_expired_at\" class=\"bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-primary-500 dark:focus:border-primary-500\" placeholder=\"Qwik\"></div><div><label for=\"status\" class=\"block mb-2 text-sm font-medium text-gray-900 dark:text-white\">Status</label> <select id=\"status\" name=\"status\" class=\"bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-primary-500 focus:border-primary-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-primary-500 dark:focus:border-primary-500\"><option value=\"\" selected disabled>Select Status</option> <option value=\"ACTIVE\">Active</option> <option value=\"NOT_ACTIVE\">Not Active</option></select></div><div class=\"\"><label for=\"commands\" class=\"block mb-2 text-sm font-medium text-gray-900 dark:text-white\">Command</label> <textarea id=\"commands\" name=\"commands\" rows=\"3\" class=\"block p-2.5 w-full text-sm text-gray-900 bg-gray-50 rounded-lg border border-gray-300 focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-primary-500 dark:focus:border-primary-500\" placeholder=\"commands\"></textarea></div><div class=\"sm:col-span-2\"><label for=\"notes\" class=\"block mb-2 text-sm font-medium text-gray-900 dark:text-white\">Notes</label> <textarea id=\"notes\" name=\"notes\" rows=\"3\" class=\"block p-2.5 w-full text-sm text-gray-900 bg-gray-50 rounded-lg border border-gray-300 focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-primary-500 dark:focus:border-primary-500\" placeholder=\"notes\"></textarea></div></div><button type=\"submit\" class=\"mt-6 inline-flex items-center justify-center px-4 py-2 text-sm font-medium tracking-wide text-white transition-colors duration-200 bg-blue-600 rounded-md hover:bg-blue-700 focus:ring-2 focus:ring-offset-2 focus:ring-blue-700 focus:shadow-outline focus:outline-none\">Create project</button></form></div></section>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
