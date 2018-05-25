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
            let tmpChartData = [];
            let chartData = {"inbound":[],"outbound":[],"count":[]};
            let tableData = [];
            for(let server in data){
                if(data[server] != null){
                    let calcData = calcServerGlobalTraffic(data[server]);
                    for(let idx in calcData.chart){
                        let tick = calcData.chart[idx];
                        if(typeof(tmpChartData[tick.time]) == "undefined"){
                            tmpChartData[tick.time] = tick;
                        } else {
                            tmpChartData[tick.time].inbound += tick.inbound;
                            tmpChartData[tick.time].outbound += tick.outbound;
                            tmpChartData[tick.time].count += tick.count;
                        }
                    }
                    console.log(tmpChartData);
                    tableData.concat(calcData.table);
                }
            }
            for(let idx in tmpChartData){
                let tick = tmpChartData[idx];
                chartData.inbound.push([tick.time,tick.inbound]);
                chartData.outbound.push([tick.time,tick.outbound]);
                chartData.count.push([tick.time,tick.count]);
            }

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
            console.log(tableData);
            genTopTable(tableData);
        }
    });
}

function calcServerGlobalTraffic(data){
    chartData = [];
    tableData = [];
    for(let idx in data){
        let tick = data[idx];
        let date = new Date(tick.time);
        let time = parseInt(date.getTime()/1000) * 1000;//去整數
        tableData.push({"date":tick.time,"time":time,"data":tick});
        let chartTick = {"time":time,"inbound":tick.traffic.inbound,"outbound":tick.traffic.outbound,"count":tick.traffic.count};
        chartData.push(chartTick);
    }
    return {"chart":chartData,"table":tableData};
}

function genTopTable(tableData) {
    //tableData.sort.sort(function(a, b){return (a.inbound>b.inbound)?1:(a.inbound==b.inbound)?-1:0});
    let tableView = '<table id="topTable" class="table"><thead><tr><th>Host/IP</th><th>Inbound</th><th>Outbound</th><th>Count</th></tr></thead><tbody>';
    for(let idx in tableData){
        let host = tableData[idx];
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
