package gaet

const testResultTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Go Appengine Test Framework Test Results</title>
    <style type="text/css" media="screen">
    body {
        background-color:#CFCFCF;
    }
    th {
        background-color:#FFEECC;
    }
    tr.fail {
        background-color:#FF0000;
    }
    tr.pass {
        background-color:#00FF00;
    }
    </style>
    <style type="text/css" media="all">
    table, td, th{
        border-width:1px;
        border-color:#000000;
        border-style:outset;
    }
    </style>
</head>
<body>
    GAET - Test Results for <b>{{.AppName}}</b>
    <table>
      <tr>
        <td>Total Number of Tests</td>
        <td>{{.TestCount}}</td>
      </tr>
      <tr>
        <td>Total Number of Passes:</td>
        <td>{{.TestPass}}</td>
      </tr>
      <tr>
        <td>Total Number of Failures:</td>
        <td>{{.TestFail}}</td>
      </tr>
      <tr>
        <td>Testing Started at:</td>
        <td>{{.Timestamp}}</td>
      </tr>
      <tr>
        <td>Total Time To Test:</td>
        <td>{{.TotalTime}} ms</td>
      </tr>
    </table>
    <br/>
    <table>
     <tr>
       <th>Test Name</th>
       <th>Test Result</th>
       <th>Test Time</th>
       <th>Test Output</th>
     </tr>
    {{range .Results}}
     <tr class="{{.TestResult}}">
       <td>{{.TestName}}</td>
       <td>{{.TestResult}}</td>
       <td>{{.TestTime}} ms</td>
       <td>{{.TestOutput}}</td>
     </tr>
     {{end}}
</body>
</html>`
