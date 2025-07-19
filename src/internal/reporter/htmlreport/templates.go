package htmlreport

import (
	"html"
	"html/template"
	"strings"
)

const summaryPageLayoutTemplate = `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
<meta http-equiv="X-UA-Compatible" content="IE=EDGE,chrome=1" />
<link href="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAMAAABEpIrGAAAAn1BMVEUAAADCAAAAAAA3yDfUAAA3yDfUAAA8PDzr6+sAAAD4+Pg3yDeQkJDTAADt7e3V1dU3yDdCQkIAAADbMTHUAABBykHUAAA2yDY3yDfr6+vTAAB3diDR0dGYcHDUAAAjhiPSAAA3yDeuAADUAAA3yDf////OCALg9+BLzktBuzRelimzKgv87+/dNTVflSn1/PWz6rO126g5yDlYniy0KgwjJ0TyAAAAI3RSTlMABAj0WD6rJcsN7X1HzMqUJyYW+/X08+bltqSeaVRBOy0cE+citBEAAADBSURBVDjLlczXEoIwFIThJPYGiL0XiL3r+z+bBOJs9JDMuLffP8v+Gxfc6aIyDQVjQcnqnvRDEQwLJYtXpZT+YhDHKIjLbS+OUeT4TjkKi6OwOArq+yeKXD9uDqQQbcOjyCy0e6bTojZSftX+U6zUQ7OuittDu1k0WHqRFfdXQijgjKfF6ZwAikvmKD6OQjmKWUcDigkztm5FZN05nMON9ZcoinlBmTNnAUdBnRbUUbgdBZwWbkcBpwXcVsBtxfjb31j1QB5qeebOAAAAAElFTkSuQmCC" rel="icon" type="image/x-icon" />
<title>{{.ReportTitle}} - {{.Translations.CoverageReport}}</title>
<link rel="stylesheet" type="text/css" href="report.css" />
<link rel="stylesheet" type="text/css" href="chartist.min.css"/>
<link rel="stylesheet" type="text/css" href="{{.AngularCssFile}}">
</head>
<body>
    <!-- Data for Angular components -->
    <script>
        window.assemblies = {{.AssembliesJSON}}; 
        window.riskHotspots = {{.RiskHotspotsJSON}}; 
        window.metrics = {{.MetricsJSON}}; 
        window.riskHotspotMetrics = {{.RiskHotspotMetricsJSON}}; 
        window.historicCoverageExecutionTimes = {{.HistoricCoverageExecutionTimesJSON}}; 
        window.translations = {{.TranslationsJSON}}; 

        window.branchCoverageAvailable = {{.BranchCoverageAvailable}};
        window.methodCoverageAvailable = {{.MethodCoverageAvailable}};
        window.maximumDecimalPlacesForCoverageQuotas = {{.MaximumDecimalPlacesForCoverageQuotas}};
    </script>

    <div class="container">
        <div class="containerleft">
            <h1>{{.ReportTitle}}
                <!-- GitHub Buttons (from C# original) -->
                <a class="button" href="https://github.com/danielpalme/ReportGenerator" title="{{.Translations.StarTooltip}}"><i class="icon-star"></i>{{.Translations.Star}}</a>
                <a class="button" href="https://github.com/sponsors/danielpalme" title="{{.Translations.SponsorTooltip}}"><i class="icon-sponsor"></i>{{.Translations.Sponsor}}</a>
            </h1>
            
            <!-- Summary Cards -->
            <div class="card-group">
                {{range .SummaryCards}}
                <div class="card">
                    <div class="card-header">{{.Title}}</div>
                    <div class="card-body">
                        {{if .ProRequired}}
                        <div class="center">
                            <p>{{$.Translations.MethodCoverageProVersion}}</p>
                            <a class="pro-button" href="https://reportgenerator.io/pro" target="_blank">{{$.Translations.MethodCoverageProButton}}</a>
                        </div>
                        {{else}}
                            {{if .SubTitle}}
                            <div class="large cardpercentagebar cardpercentagebar{{.SubTitlePercentageBarValue}}">{{.SubTitle}}</div>
                            {{end}}
                            <div class="table">
                                <table>
                                    {{range .Rows}}
                                    <tr><th>{{.Header}}:</th><td class="limit-width {{if eq .Alignment "right"}}right{{end}}" title="{{.Tooltip}}">{{.Text}}</td></tr>
                                    {{end}}
                                </table>
                            </div>
                        {{end}}
                    </div>
                </div>
                {{end}}
            </div>

            <!-- Overall History Chart -->
            {{if .OverallHistoryChartData.Series}}
                <h1>{{.Translations.History}}</h1>
                <!-- The SVG is rendered directly by Go, Chartist.js might not be needed for this if SVG is static -->
                <div class="historychart ct-chart" data-data="historyChartDataOverall">{{.OverallHistoryChartData.SVGContent | SafeHTML}}</div>
                <!-- If custom.js or Angular needs the data for interactivity with this chart: -->
                <!-- <script type="text/javascript">/* <![CDATA[ */ 
                // var historyChartDataOverall = {{.OverallHistoryChartData.JSONData | SafeJS}};
                // /* ]]> */ </script> -->
            {{end}}

            <!-- Risk Hotspots Section (Angular Component) -->
            <h1>{{.Translations.RiskHotspots}}</h1>
            <risk-hotspots></risk-hotspots> 
            {{if not .HasRiskHotspots}}
            <p>{{.Translations.NoRiskHotspots}}</p>
            {{end}}

            <!-- Coverage Section (Angular Component) -->
            <h1>{{.Translations.Coverage3}}</h1>
            <coverage-info></coverage-info> 
            {{if not .HasAssemblies}}
            <p>{{.Translations.NoCoveredAssemblies}}</p>
            {{end}}

            <div class="footer">{{.Translations.GeneratedBy}} ReportGenerator {{.AppVersion}}<br />{{.CurrentDateTime}}<br /><a href="https://github.com/danielpalme/ReportGenerator">GitHub</a> | <a href="https://reportgenerator.io">reportgenerator.io</a></div>
        </div> <!-- End containerleft -->
    </div> <!-- End container -->

    <script type="text/javascript" src="chartist.min.js"></script> <!-- For Angular components if they use Chartist -->
    <script type="text/javascript" src="custom.js"></script>
    <script type="text/javascript" src="{{.CombinedAngularJsFile}}"></script>
</body>
</html>`

const classDetailLayoutTemplate = `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
<meta http-equiv="X-UA-Compatible" content="IE=EDGE,chrome=1" />
<link href="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAMAAABEpIrGAAAAn1BMVEUAAADCAAAAAAA3yDfUAAA3yDfUAAA8PDzr6+sAAAD4+Pg3yDeQkJDTAADt7e3V1dU3yDdCQkIAAADbMTHUAABBykHUAAA2yDY3yDfr6+vTAAB3diDR0dGYcHDUAAAjhiPSAAA3yDeuAADUAAA3yDf////OCALg9+BLzktBuzRelimzKgv87+/dNTVflSn1/PWz6rO126g5yDlYniy0KgwjJ0TyAAAAI3RSTlMABAj0WD6rJcsN7X1HzMqUJyYW+/X08+bltqSeaVRBOy0cE+citBEAAADBSURBVDjLlczXEoIwFIThJPYGiL0XiL3r+z+bBOJs9JDMuLffP8v+Gxfc6aIyDQVjQcnqnvRDEQwLJYtXpZT+YhDHKIjLbS+OUeT4TjkKi6OwOArq+yeKXD9uDqQQbcOjyCy0e6bTojZSftX+U6zUQ7OuittDu1k0WHqRFfdXQijgjKfF6ZwAikvmKD6OQjmKWUcDigkztm5FZN05nMON9ZcoinlBmTNnAUdBnRbUUbgdBZwWbkcBpwXcVsBtxfjb31j1QB5qeebOAAAAAElFTkSuQmCC" rel="icon" type="image/x-icon" />
<title>{{.Class.Name}} - {{.ReportTitle}}</title>
<link rel="stylesheet" type="text/css" href="report.css" />
<link rel="stylesheet" type="text/css" href="{{.AngularCssFile}}">
</head>
<body>
    <script>
        window.classDetails = JSON.parse({{.ClassDetailJSON}});
        window.assemblies = JSON.parse({{.AssembliesJSON}});
        window.translations = JSON.parse({{.TranslationsJSON}});
        window.branchCoverageAvailable = {{.BranchCoverageAvailable}};
        window.methodCoverageAvailable = {{.MethodCoverageAvailable}};
        window.maximumDecimalPlacesForCoverageQuotas = {{.MaximumDecimalPlacesForCoverageQuotas}};
        window.riskHotspots = JSON.parse({{.RiskHotspotsJSON}}); 
        window.metrics = JSON.parse({{.MetricsJSON}});
        window.riskHotspotMetrics = JSON.parse({{.RiskHotspotMetricsJSON}});
        window.historicCoverageExecutionTimes = JSON.parse({{.HistoricCoverageExecutionTimesJSON}});
    </script>

    <div class="container">
        <div class="containerleft">
            <h1><a href="index.html" class="back"><</a> {{.Translations.Summary}}</h1>

            <div class="card-group">
                <div class="card">
                    <div class="card-header">{{.Translations.Information}}</div>
                    <div class="card-body">
                        <div class="table">
                            <table>
                                <tr><th>{{.Translations.Class}}:</th><td class="limit-width" title="{{.Class.Name}}">{{.Class.Name}}</td></tr>
                                <tr><th>{{.Translations.Assembly}}:</th><td class="limit-width" title="{{.Class.AssemblyName}}">{{.Class.AssemblyName}}</td></tr>
                                <tr><th>{{.Translations.Files3}}:</th><td class="overflow-wrap">
                                    {{$filesLen := len .Class.Files}}
                                    {{$lastFileIdx := sub $filesLen 1}}
                                    {{range $idx, $file := .Class.Files}}
                                        <a href="#{{$file.ShortPath}}" class="navigatetohash">{{$.Translations.File}} {{$idx | inc}}: {{$file.Path}}</a>{{if ne $idx $lastFileIdx}}<br />{{end}}
                                    {{else}}
                                        No files found.
                                    {{end}}
                                </td></tr>
                                {{if .Tag}}
                                <tr><th>{{.Translations.Tag}}:</th><td class="limit-width" title="{{.Tag}}">{{.Tag}}</td></tr>
                                {{end}}
                            </table>
                        </div>
                    </div>
                </div>
            </div>

            <div class="card-group">
                <div class="card">
                    <div class="card-header">{{.Translations.LineCoverage}}</div>
                    <div class="card-body">
                        <div class="large cardpercentagebar cardpercentagebar{{.Class.CoveragePercentageBarValue}}">{{.Class.CoveragePercentageForDisplay}}</div>
                        <div class="table">
                            <table>
                                <tr><th>{{.Translations.CoveredLines}}:</th><td class="limit-width right" title="{{.Class.CoveredLines}}">{{.Class.CoveredLines}}</td></tr>
                                <tr><th>{{.Translations.UncoveredLines}}:</th><td class="limit-width right" title="{{.Class.UncoveredLines}}">{{.Class.UncoveredLines}}</td></tr>
                                <tr><th>{{.Translations.CoverableLines}}:</th><td class="limit-width right" title="{{.Class.CoverableLines}}">{{.Class.CoverableLines}}</td></tr>
                                <tr><th>{{.Translations.TotalLines}}:</th><td class="limit-width right" title="{{.Class.TotalLines}}">{{.Class.TotalLines}}</td></tr>
                                <tr><th>{{.Translations.LineCoverage}}:</th><td class="limit-width right" title="{{.Class.CoveredLines}} of {{.Class.CoverableLines}}">{{.Class.CoverageRatioTextForDisplay}}</td></tr>
                            </table>
                        </div>
                    </div>
                </div>
                {{if .BranchCoverageAvailable}}
                <div class="card">
                    <div class="card-header">{{.Translations.BranchCoverage}}</div>
                    <div class="card-body">
                        <div class="large cardpercentagebar cardpercentagebar{{.Class.BranchCoveragePercentageBarValue}}">{{.Class.BranchCoveragePercentageForDisplay}}</div>
                        <div class="table">
                            <table>
                                <tr><th>{{.Translations.CoveredBranches2}}:</th><td class="limit-width right" title="{{.Class.CoveredBranches}}">{{.Class.CoveredBranches}}</td></tr>
                                <tr><th>{{.Translations.TotalBranches}}:</th><td class="limit-width right" title="{{.Class.TotalBranches}}">{{.Class.TotalBranches}}</td></tr>
                                <tr><th>{{.Translations.BranchCoverage}}:</th><td class="limit-width right" title="{{.Class.CoveredBranches}} of {{.Class.TotalBranches}}">{{.Class.BranchCoverageRatioTextForDisplay}}</td></tr>
                            </table>
                        </div>
                    </div>
                </div>
                {{end}}
                 <div class="card">
                    <div class="card-header">{{.Translations.MethodCoverage}}</div>
                    <div class="card-body">
                        {{if .MethodCoverageAvailable}}
                        <div class="large cardpercentagebar cardpercentagebar{{.Class.MethodCoveragePercentageBarValue}}">{{.Class.MethodCoveragePercentageForDisplay}}</div>
                        <div class="table">
                            <table>
                                <tr><th>{{.Translations.CoveredCodeElements}}:</th><td class="limit-width right" title="{{.Class.CoveredMethods}}">{{.Class.CoveredMethods}}</td></tr>
                                <tr><th>{{.Translations.FullCoveredCodeElements}}:</th><td class="limit-width right" title="{{.Class.FullyCoveredMethods}}">{{.Class.FullyCoveredMethods}}</td></tr>
                                <tr><th>{{.Translations.TotalCodeElements}}:</th><td class="limit-width right" title="{{.Class.TotalMethods}}">{{.Class.TotalMethods}}</td></tr>
                                <tr><th>{{.Translations.CodeElementCoverageQuota2}}:</th><td class="limit-width right" title="{{.Class.CoveredMethods}} of {{.Class.TotalMethods}}">{{.Class.MethodCoverageRatioTextForDisplay}}</td></tr>
                                <tr><th>{{.Translations.FullCodeElementCoverageQuota2}}:</th><td class="limit-width right" title="{{.Class.FullyCoveredMethods}} of {{.Class.TotalMethods}}">{{.Class.FullMethodCoverageRatioTextForDisplay}}</td></tr>
                            </table>
                        </div>
                        {{else}}
                        <div class="center">
                            <p>{{.Translations.MethodCoverageProVersion}}</p>
                            <a class="pro-button" href="https://reportgenerator.io/pro" target="_blank">{{.Translations.MethodCoverageProButton}}</a>
                        </div>
                        {{end}}
                    </div>
                </div>
            </div>

            {{if .Class.FilesWithMetrics}} <!-- This condition might need to be based on .Class.MetricsTable.Rows now -->
            <h1>{{.Translations.Metrics}}</h1>
            <div class="table-responsive">
                <table class="overview table-fixed">
                    <colgroup>
                        <col class="column-min-200" />
                        {{range .Class.MetricsTable.Headers}}
                        <col class="column105" />
                        {{end}}
                    </colgroup>
                    <thead><tr><th>{{$.Translations.Methods}}</th>
                        {{range .Class.MetricsTable.Headers}}
                        <th>{{.Name}} {{if .ExplanationURL}}<a href="{{.ExplanationURL}}" target="_blank"><i class="icon-info-circled"></i></a>{{end}}</th>
                        {{end}}
                    </tr></thead>
                    <tbody>
                        {{range .Class.MetricsTable.Rows}}
                        <tr><td title="{{.FullName}}"><a href="#{{.FileShortPath}}_line{{.Line}}" class="navigatetohash">{{if $.Class.IsMultiFile}}File {{.FileIndexPlus1}}: {{end}}{{.Name}}</a></td>
                            {{range .MetricValues}}<td>{{.}}</td>{{end}}
                        </tr>
                        {{end}}
                    </tbody>
                </table>
            </div>
            {{end}}

            <h1>{{.Translations.Files3}}</h1>
            {{range $fileIdx, $file := .Class.Files}}
            <h2 id="{{$file.ShortPath}}">{{$file.Path}}</h2>
            <div class="table-responsive">
                <table class="lineAnalysis">
                    <thead><tr><th></th><th>#</th><th>{{$.Translations.Line}}</th><th></th><th>{{$.Translations.LineCoverage}}</th></tr></thead>
                    <tbody>
                    {{range $file.Lines}}
                        <tr class="{{if ne .LineVisitStatus "gray"}}coverableline{{end}}" title="{{.Tooltip}}" data-coverage="{{.DataCoverage}}">
                            <td class="{{.LineVisitStatus}}"> </td>
                            <td class="leftmargin rightmargin right">{{if ne .LineVisitStatus "gray"}}{{.Hits}}{{end}}</td>
                            <td class="rightmargin right"><a id="{{$file.ShortPath}}_line{{.LineNumber}}"></a><code>{{.LineNumber}}</code></td>
                            {{if .IsBranch}}
                            <td class="percentagebar percentagebar{{.BranchBarValue}}"><i class="icon-fork"></i></td>
                            {{else}}
                            <td></td>
                            {{end}}
                            <td class="light{{.LineVisitStatus}}"><code>{{.LineContent | SanitizeSourceLine}}</code></td>
                        </tr>
                    {{end}}
                    </tbody>
                </table>
            </div>
            {{else}}
                <p>{{.Translations.NoFilesFound}}</p>
            {{end}}

            <div class="footer">{{.Translations.GeneratedBy}} ReportGenerator {{.AppVersion}}<br />{{.CurrentDateTime}}<br /><a href="https://github.com/danielpalme/ReportGenerator">GitHub</a> | <a href="https://reportgenerator.io">reportgenerator.io</a></div>
        </div> 

        {{if .Class.SidebarElements}}
        <div class="containerright">
            <div class="containerrightfixed">
                <h1>{{.Translations.MethodsProperties}}</h1>
                {{range .Class.SidebarElements}}
                <a href="#{{.FileShortPath}}_line{{.Line}}" class="navigatetohash percentagebar percentagebar{{.CoverageBarValue}}" title="{{if $.Class.IsMultiFile}}File {{.FileIndexPlus1}}: {{end}}{{.CoverageTitle}} - {{.Name}}"><i class="icon-{{.Icon}}"></i>{{.Name}}</a><br />
                {{end}}
                <br/>
            </div>
        </div>
        {{end}}
    </div> 

    <script type="text/javascript" src="custom.js"></script> 
    <script type="text/javascript" src="{{.CombinedAngularJsFile}}"></script>
</body>
</html>`

var (
	// classDetailTpl for class detail pages (server-rendered structure)
	classDetailTpl = template.Must(template.New("classDetail").Funcs(template.FuncMap{
		"inc":      func(i int) int { return i + 1 },
		"sub":      func(a, b int) int { return a - b },
		"SafeHTML": func(s string) template.HTML { return template.HTML(s) },
		"SafeJS":   func(s string) template.JS { return template.JS(s) },
		"SanitizeSourceLine": func(line string) template.HTML {
			// 1. HTML-escape first to get &lt;, &gt;, &amp; …
			escaped := html.EscapeString(line)

			// 2. Replace TABs with four real spaces first (so that step 3 sees them)
			escaped = strings.ReplaceAll(escaped, "\t", "    ")

			// 3. Turn every real space into &nbsp;
			escaped = strings.ReplaceAll(escaped, " ", "&nbsp;")

			return template.HTML(escaped) // mark it safe – we built the HTML ourselves
		},
	}).Parse(classDetailLayoutTemplate))

	// summaryPageTpl for the main index.html (summary page)
	summaryPageTpl = template.Must(template.New("summaryPage").Funcs(template.FuncMap{
		"SafeHTML": func(s string) template.HTML { return template.HTML(s) },
		"SafeJS":   func(s string) template.JS { return template.JS(s) },
	}).Parse(summaryPageLayoutTemplate))
)
