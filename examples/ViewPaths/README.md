## Demonstrates the basic setup of a site with an URL structure

Views are defined in sub-directories of the "views" directory.
The directory structure mimics the URL structure with the exception
that root views are defined in the "views/root" directory instead of in "views".

"views" contains a paths.go file where variables for all views and the
URL structure get declared. "views" has no dependencies on its sub packages,
so views with can be imported in all sub packages to access the variables
for all views (usually to get their URL).

If views have URL arguments, then the index of the URL argument is added
to the package name (see user_0 below).

* Project directory
	* config.json
	* main.go
	* views
		* paths.go
		* root
			* homepage.go
			* getjson.go
			* getxml.go
		* admin
			* admin.go
			* user_0
				* user_0.go


Download, build and run example:

	go get github.com/ungerik/go-start/examples/ViewPaths
	go install github.com/ungerik/go-start/examples/ViewPaths && ViewPaths
