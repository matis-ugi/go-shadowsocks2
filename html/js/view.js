function showManagementView() {
    Highcharts.setOptions({
    global: {
        useUTC: false
    }
});
    $("#login").fadeOut("slow",function(){
        let manView = '<div id="chart" class="chart"><div id="main_chart"></div><div id="top" class="top"></div></div>';
        $("#view").empty();
        $("#view").append(manView);
        $("#view").fadeIn("slow");
        getData();
    });
}

function getData(){
    getChartData();
}

function getChartData() {
    end = parseInt(new Date().getTime()/1000);
    start = end - 60 * 60;
    $.ajax({
        url: 'api/global_chart/global/'+start+'/'+end,
        type: 'POST',
        data: JSON.stringify(User),
        contentType: 'application/json; charset=utf-8',
        dataType: 'json',
        async: false,
        success: function(data) {
            chartData = {"date":[],"inbound":[],"outbound":[],"count":[],"inbound_diff":[],"outbound_diff":[],"count_diff":[]};
            tableData = {"sort":[],"raw":[]};
            for(let server in data){
                let calcData = calcServerGlobalTraffic(data[server]);
                chartData = calcData[0];
                /*for(let idx in calcData[0]){
                    let tick = calcData[0][idx];
                    chartData.inbound.push([time,tick.inbound]);
                    chartData.outbound.push([time,tick.outbound]);
                    chartData.count.push([time,tick.count]);
                }*/
            }
            console.log(chartData);
            /*chartCalcData = {"inbound":[],"outbound":[],"count":[],"inbound_diff":[],"outbound_diff":[],"count_diff":[]};
            for(let idx in chartData.date){
                let time = chartData.date[idx];
                let inbound = [time,chartData.inbound_diff[idx]];
                chartCalcData.inbound_diff.push(inbound);
                let outbound = [time,chartData.outbound_diff[idx]];
                chartCalcData.outbound_diff.push(outbound);
                let count = [time,chartData.count_diff[idx]];
                chartCalcData.count_diff.push(count);
                 inbound = [time,chartData.inbound[idx]];
                chartCalcData.inbound.push(inbound);
                 outbound = [time,chartData.outbound[idx]];
                chartCalcData.outbound.push(outbound);
                 count = [time,chartData.count[idx]];
                chartCalcData.count.push(count);
            }*/
            genChart(chartData);
            genTopTable(tableData);
        }
    });
}

function calcServerGlobalTraffic(data){
    chartData = {"date":[],"inbound":[],"outbound":[],"count":[],"inbound_diff":[],"outbound_diff":[],"count_diff":[]};
    tableData = {"sort":[],"raw":[]};
    for(let idx in data){
        let tick = data[idx];
        let date = new Date(tick.time);
        tableData.raw.push({"date":tick.time,"time":date.getTime(),"data":tick});
        let time = date.getTime();
        chartData.inbound.push([time,tick.traffic.inbound]);
        chartData.outbound.push([time,tick.traffic.outbound]);
        chartData.count.push([time,tick.traffic.count]);
    }
    return [chartData,tableData];
}

function genTopTable(tableData) {
    //tableData.sort.sort(function(a, b){return (a.inbound>b.inbound)?1:(a.inbound==b.inbound)?-1:0});
    let tableView = '<table id="topTable" class="table"><thead><tr><th>Host/IP</th><th>Inbound</th><th>Outbound</th><th>Count</th></tr></thead><tbody>';
    for(let idx in tableData.sort){
        let host = tableData.sort[idx];
        tableView += '<tr><td>'+idx+'</td><td>'+host.inbound+'</td><td>'+host.outbound+'</td><td>'+host.count+'</td></tr>';
    }
    tableView += '</tbody></table>';
    $("#top").empty();
    $("#top").append(tableView);
    $('#topTable').DataTable({order: [[ 0, 'desc' ], [ 0, 'asc' ]]});
    setTimeout(function(){

    },2000);

}

function genChart(data){
    Highcharts.stockChart('main_chart', {

    title: {
        text: "Traffic",
    },
    tooltip: {
        crosshairs: true,
        shared: true
    },
    plotOptions: {
        series: {
            label: {
                connectorAllowed: true
            },
             pointStart: 0
        }
    },

    series: [
        {
            name: 'Inbound',
            data: data.inbound,
        },
        {
            name: 'Outbound',
            data: data.outbound,
        },
        {
            name: 'Count',
            data: data.count,
        }
    ],

    responsive: {
        rules: [{
            condition: {
                maxWidth: 500
            },
            chartOptions: {
                legend: {
                    layout: 'horizontal',
                    align: 'center',
                    verticalAlign: 'bottom'
                }
            }
        }]
    }

});
}
