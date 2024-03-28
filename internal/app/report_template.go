package app

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Report {{ .Year }} week {{ .WeekNumber }}</title>
	<style>
		* {	
			font-family: sans-serif;
		}
	</style>
</head>

<body>
	<h1>Report {{ .Year }} week {{ .WeekNumber }}</h1>
	<ul>
{{- range $task := .Tasks }}
	<li><a href={{ $task.Link }}>{{ $task.Key }}</a> | {{ $task.Name }} - <strong>{{ $task.Status }}</strong></li>
{{- end }}
	</ul>
</body>
</html>
`
