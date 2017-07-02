"use strict";

(function(){
    var d = document,
        tradeHistory = null,
        lastDateFilter = "All Dates",
        dateFilter = "All Dates",
        lastStart = -1,
        lastEnd = -1,
        orderCache = {},
        valueCache = {},
        initialOrderBook = $(".orderBook tbody"),
        zero = BigNumber(0);
    
    $("#analyzeBtn").removeAttr('disabled');
    $("#clearBtn").click(clearAnalysis);

    $('.allIndicator').on('click', function(e) {
        var node = $(this),
            orderBook = node.closest('table').find('tbody'),
            children,
            childNode,
            i, ii, method;

        if (!node.closest('.analysis').length) {
            return; // Do nothing is not in analysis mode
        }

        node.removeClass('mixed');

        if (node.hasClass('all')) {
            node.removeClass('all');
            method = removeOrder;
            children = orderBook.children().filter('.selected');
            valueCache.selected = 0;
        } else {
            node.addClass('all');
            method = addOrder;
            children = orderBook.children().not('.selected');
            valueCache.selected = valueCache.totalOrders;
        }
        
        for (i = 0, ii = children.length; i < ii; ++i) {
            childNode = $(children[i]);
            method(orderCache[childNode.attr('data-id')]);
            childNode.toggleClass('selected');
        }

        refreshAnalysis();
    });

    $("#dateWidget").datepick({
        minDate: moment.utc('2014-01-01').toDate(),
        maxDate: moment.utc(new Date).endOf('month').toDate(),
        rangeSelect: true,
        monthsToShow: 2,
        onSelect: function() {
            dateFilter = "User Defined";
            applyDateRange();
        }
    });

    $("#dateRangeBtn").click(function() {
        $("#dateRangeSelector").show();
    });

    $("#closeDateRange").click(function() {
        $("#dateRangeSelector").fadeOut();
        dateFilter = lastDateFilter;
        applyDateRange();
    });
		
		$("#coinName").on("keyup",function(){
			if ($("#coinName").val().length>0){
				$("#analyzeBtn").removeAttr('disabled');
			} else {
				$("#analyzeBtn").attr('disabled', true);
			}
		});
		
    $("#coinName").autocomplete({
        source: function(req, res) {
            var r = $.ui.autocomplete.filter(Object.keys(CURRENCIES), req.term);

            res(r.slice(0, 10));
        },
        select: function(evt, ui) {
            evt.preventDefault();
           
            // Set to actual token value, not visible name
            $(this).val(CURRENCIES[ui.item.value].token);
        }
    });

    function refreshAnalyzeControls() {
        if ($("#coinName").val().trim() == "" ||
            $("#marketName").val().trim() == "") {

            $("#analyzeBtn").attr('disabled', true);
            $("#showList, #dateRangeBtn").removeAttr('disabled');
        } else {

            $("#analyzeBtn, #showList, #dateRangeBtn").removeAttr('disabled');
        }
    }

    $("#coinName, #marketName").change(refreshAnalyzeControls);

    function toPHPDate(dt) {
        return +dt / 1000 | 0;
    }

    function getDateRange() {
        var start, end, dates;
        
        if (dateFilter == "Today") {
            start = toPHPDate(moment.utc().startOf('day').toDate());
            end = toPHPDate(moment.utc().endOf('day').toDate());
        } else if (dateFilter == "This Week") {
            start = toPHPDate(moment.utc().startOf('week').toDate());
            end = toPHPDate(moment.utc().endOf('week').toDate());
        } else if (dateFilter == "This Month") {
            start = toPHPDate(moment.utc().startOf('month').toDate());
            end = toPHPDate(moment.utc().endOf('month').toDate());
        } else if (dateFilter == "This Year") {
            start = toPHPDate(moment.utc().startOf('year').toDate());
            end = toPHPDate(moment.utc().endOf('year').toDate());
        } else if (dateFilter == "Yesterday") {
            start = toPHPDate(moment.utc().subtract('days', 1).startOf('day').toDate());
            end = toPHPDate(moment.utc().subtract('days', 1).endOf('day').toDate());
        } else if (dateFilter == "Last Week") {
            start = toPHPDate(moment.utc().subtract('weeks', 1).startOf('week').toDate());
            end = toPHPDate(moment.utc().subtract('weeks', 1).endOf('week').toDate());
        } else if (dateFilter == "Last Month") {
            start = toPHPDate(moment.utc().subtract('months', 1).startOf('month').toDate());
            end = toPHPDate(moment.utc().subtract('months', 1).endOf('month').toDate());
        } else if (dateFilter == "All Dates") {
            start = 0;
            end = 1917895987;
        } else if (dateFilter == "User Defined") {
            dates = $("#dateWidget").datepick('getDate');
            start = toPHPDate(startOfDate(dates[0]));
            end = toPHPDate(endOfDate(dates[1]));
        }

        return {start: start, end: end};
    }

    function applyDateRange() {
        var range, start, end, startFmt, endFmt
        
        range = getDateRange();
        start = range.start;
        end = range.end;

        if (dateFilter != "User Defined") {
            $(".dateRange").text('From ' + dateFilter);
            
            if (start == 0) {
                setWidgetDate($("#dateWidget"), [new Date, new Date]);
            } else {
                setWidgetDate($("#dateWidget"), [
                    new Date(start * 1000),
                    new Date(end * 1000)
                ]);
            }
        } else {
            startFmt = moment.utc(new Date(start * 1000)).format("YYYY-MM-DD");
            endFmt = moment.utc(new Date(end * 1000)).format("YYYY-MM-DD");
            $(".dateRange").text(startFmt + " to " + endFmt);
        }

        //updateHash();
        //$("#dateRangeSelector").fadeOut();
    };

    $(".rangeLinks a").click(function(e) {
        e.preventDefault();
        dateFilter = $(this).text().trim();
        
        applyDateRange();
    });

    function parseHash(data) {
        var o;

        if (!/^#analysis\//.test(data)) {
            return null;
        }
        
        data = data.substring(10).split('/');

        o = {
            coin: data[0],
            market: data[1],
            show: data[2],
            fromDate: parseInt(data[3], 10),
            toDate: parseInt(data[4], 10)
        }

        if (o.fromDate !== o.fromDate) { // NaN check
            delete o.fromDate;
        }
        
        if (o.toDate !== o.toDate) { // NaN check
            delete o.toDate;
        }

        return o;
    }

    function clearAnalysis(e) {
        e.preventDefault();

        valueCache = {};

        $(".orderBook").removeClass('analysis');
        $(".tradeAnalysis").removeClass('active');

        $(".orderBook tbody").replaceWith(initialOrderBook);
        
        $("#clearBtn").attr('disabled', true);
        
        $(".calc").hide();
        $(".info").show();
        $(".pages").show();
        $("#coinName").val('');
        dateFilter = "All Dates";
        $(".dateRange").text('From ' + dateFilter);
        applyDateRange();

        location.hash = '';
        updateHash();
    };

    function refreshAnalysis() {
        var breakEvenPrice = zero,
            total, totalAmount,
            buyFeeAvg, oneMinusBuyFee,
            totalFees, profitLoss,
            totalPlusFees;

		$('.profitLoss').removeClass("sellClass buyClass");
		$("#tradeHistoryPages").empty();
		analysis=true;
		
        if (valueCache.buyAmountTotal.equals(zero)) {
            $(".avgBuyPrice, .totalBuys").text('--');
        } else {
            $(".avgBuyPrice").text(
                valueCache.buyTotal.plus(valueCache.buyLendingTotal).dividedBy(valueCache.buyAmountTotal).toFixed(8) + " " + valueCache.market);
            $(".totalBuys").text(valueCache.buyAmountTotal.toFixed(8) + " " + valueCache.coin);
        }

        if (valueCache.sellAmountTotal.equals(zero)) {
            $(".avgSellPrice, .totalSells").text('--');
        } else {
            $(".avgSellPrice").text(
                valueCache.sellTotal.plus(valueCache.sellLendingTotal).dividedBy(valueCache.sellAmountTotal).toFixed(8) + " " + valueCache.market);
            $(".totalSells").text(valueCache.sellAmountTotal.toFixed(8) + " " + valueCache.coin);
        }

        if (valueCache.buyAmountTotal.equals(zero) &&
            valueCache.sellAmountTotal.equals(zero)) {

            $(".breakEvenPrice, .totalBuys, .profitLoss, .totalSells").text('--');
            return;
        }
        
        if (valueCache.market == "-"){
	        $(".breakEvenPrice, .avgBuyPrice, .totalBuys, .avgSellPrice, .totalSells").text('--');
	        $(".profitLossLabel").text("Total with Fees:");
			updateTooltip($(".profitLossLabel"), 'Total ' + valueCache.coin + ' spent or earned, with fees.');
			
            $(".profitLoss").text(valueCache.buyTotal
                                  .plus(valueCache.buyFeeValueTotal)
                                  .toFixed(8) + " " + valueCache.coin);
	        return;
        }

        total = valueCache.breakEvenTotal;
        totalAmount = valueCache.breakEvenAmountTotal;
        buyFeeAvg = valueCache.buyFeeTotal.dividedBy(valueCache.buyFeeCount);
        oneMinusBuyFee = BigNumber(1).minus(buyFeeAvg);

        if (valueCache.buyTotal.equals(zero)) {
            $(".profitLossLabel").text("Total with Fees:");
			updateTooltip($(".profitLossLabel"), 'Total BTC spent or earned, with fees.');
			
            $(".profitLoss").text(valueCache.sellTotal
                                  .plus(valueCache.sellFeeValueTotal)
                                  .minus(valueCache.sellLendingTotal)
                                  .toFixed(8) + " " + valueCache.market);
        
        } else if (valueCache.sellTotal.equals(zero)) {
            $(".profitLossLabel").text("Total with Fees:");
			updateTooltip($(".profitLossLabel"), 'Total BTC spent or earned, with fees.');
			
            $(".profitLoss").text(valueCache.buyTotal
                                  .plus(valueCache.buyFeeValueTotal)
                                  .toFixed(8) + " " + valueCache.market);
        
        } else {
            $(".profitLossLabel").text("Profit/Loss:");
			updateTooltip($(".profitLossLabel"), 'Total bought less total sold in BTC, with fees.');
			
            profitLoss = valueCache.sellTotal.minus(valueCache.buyTotal);

            $(".profitLoss").text((profitLoss.isNegative() ? "" : "+") +
                                  profitLoss.toFixed(8) + " " + 
                                  valueCache.market)
							.addClass(profitLoss.isNegative() ? 'sellClass' : 'buyClass');
        }
        
        if (totalAmount.greaterThan(0)) {
            breakEvenPrice = total.dividedBy(totalAmount.times(oneMinusBuyFee));
        } else if (totalAmount.lessThan(0)) {
            breakEvenPrice = total.times(oneMinusBuyFee).dividedBy(totalAmount);
        }
       
        if (totalAmount.equals(0) || 
            (profitLoss != null && profitLoss.gte(0)) ||
            valueCache.sellAmountTotal.gte(valueCache.buyAmountTotal)) {
            $(".breakEvenPrice").text('--');
        } else {
            $(".breakEvenPrice").text(breakEvenPrice.toFixed(8) + " " + valueCache.market);
        }

    }

    function makeTextCell(text) {
        return $(d.createElement('td')).text(text);
    }

    function makeNumberCell(text) {
        return makeTextCell(text).css('text-align', 'right');
    }

    function mangleOrder(order) {
        order.fee = BigNumber(order.fee);
        order.amount = BigNumber(order.amount);
        order.total = BigNumber(order.total);
        order.rate = BigNumber(order.rate);
        return order;
    }

    function addOrder(order) {
        order = mangleOrder(order);

        if (order.type == "buy") {
            addBuy(order);
        } else if (order.type == "sell") {
            addSell(order);
        }
    }

    function removeOrder(order) {
        order = mangleOrder(order);

        if (order.type == "buy") {
            removeBuy(order);
        } else if (order.type == "sell") {
            removeSell(order);
        }
    }

    function addBuy(order) {
        var oneMinusFee = BigNumber(1).minus(order.fee);

		if (order.rate <= 0){
			// valueCache.buyAmountTotal = valueCache.buyAmountTotal.minus(order.total);
	    } else {
	        // Avg buy
	        valueCache.buyTotal = valueCache.buyTotal.plus(order.total);
	        valueCache.buyAmountTotal = valueCache.buyAmountTotal.plus(order.amount.times(oneMinusFee));
	
	        // Break even
	        valueCache.breakEvenTotal = valueCache.breakEvenTotal.plus(order.rate.times(order.amount));
	        valueCache.breakEvenAmountTotal = valueCache.breakEvenAmountTotal.plus(
	            order.amount.times(oneMinusFee));
	
	        //valueCache.buyFeeTotal = valueCache.buyFeeTotal.plus(order.fee);
	        //valueCache.buyFeeValueTotal = valueCache.buyFeeValueTotal.plus(order.fee.times(order.amount));
	        valueCache.buyFeeCount++;
        }
    }

    function removeBuy(order) {
        var oneMinusFee = BigNumber(1).minus(order.fee);

		if (order.rate <= 0){
	        // valueCache.buyAmountTotal = valueCache.buyAmountTotal.plus(order.total);
		} else {
	        // Avg buy
	        valueCache.buyTotal = valueCache.buyTotal.minus(order.total);
	        valueCache.buyAmountTotal = valueCache.buyAmountTotal.minus(order.amount.times(oneMinusFee));
	
	        // Break even
	        valueCache.breakEvenTotal = valueCache.breakEvenTotal.minus(order.rate.times(order.amount));
	        valueCache.breakEvenAmountTotal = valueCache.breakEvenAmountTotal.minus(
	            order.amount.times(oneMinusFee));
	        
	        //valueCache.buyFeeTotal = valueCache.buyFeeTotal.minus(order.fee);
	        //valueCache.buyFeeValueTotal = valueCache.buyFeeValueTotal.minus(order.fee.times(order.amount));
	        valueCache.buyFeeCount--;
	    }
    }

    function addSell(order) {
        var oneMinusFee = BigNumber(1).minus(order.fee);

		if (order.rate <= 0){
			valueCache.sellTotal = valueCache.sellTotal.minus(order.total);
			valueCache.sellLendingTotal = valueCache.sellLendingTotal.plus(order.total);
			valueCache.breakEvenTotal = valueCache.breakEvenTotal.plus(order.total);
			
			valueCache.sellFeeValueTotal = valueCache.sellFeeValueTotal.plus(order.total);
	        valueCache.sellFeeCount++;
		} else {

	        // Avg sell
	        valueCache.sellTotal = valueCache.sellTotal.plus(order.total.times(oneMinusFee));
	        valueCache.sellAmountTotal = valueCache.sellAmountTotal.plus(order.amount);
	
	        // Break even
	        valueCache.breakEvenTotal = valueCache.breakEvenTotal.minus(
	            order.rate.times(order.amount).times(oneMinusFee));
	        valueCache.breakEvenAmountTotal = valueCache.breakEvenAmountTotal.minus(order.amount);
	
	        valueCache.sellFeeTotal = valueCache.sellFeeTotal.plus(order.fee);
	        valueCache.sellFeeValueTotal = valueCache.sellFeeValueTotal.plus(order.fee.times(order.total));
	        valueCache.sellFeeCount++;
        }
    }

    function removeSell(order) {
        var oneMinusFee = BigNumber(1).minus(order.fee);
        
        if (order.rate <= 0){
	        valueCache.sellTotal = valueCache.sellTotal.plus(order.total);
	        valueCache.sellLendingTotal = valueCache.sellLendingTotal.minus(order.total);
	        valueCache.breakEvenTotal = valueCache.breakEvenTotal.minus(order.total);
	        
			valueCache.sellFeeValueTotal = valueCache.sellFeeValueTotal.minus(order.total);
	        valueCache.sellFeeCount--;
		} else {
        // Avg sell
	        valueCache.sellTotal = valueCache.sellTotal.minus(order.total.times(oneMinusFee));
	        valueCache.sellAmountTotal = valueCache.sellAmountTotal.minus(order.amount);
	
	        // Break even
	        valueCache.breakEvenTotal = valueCache.breakEvenTotal.plus(
	            order.rate.times(order.amount).times(oneMinusFee));
	        valueCache.breakEvenAmountTotal = valueCache.breakEvenAmountTotal.plus(order.amount);
	        
	        valueCache.sellFeeTotal = valueCache.sellFeeTotal.minus(order.fee);
	        valueCache.sellFeeValueTotal = valueCache.sellFeeValueTotal.minus(order.fee.times(order.total));
	        valueCache.sellFeeCount--;
        }
    }

    function formatFee(order,feeCurrency) {
        var fee = BigNumber(order.fee),
            total = BigNumber(order.total),
            amount = BigNumber(order.amount), totalFee;
        if (order.type=="sell"){
          totalFee = total.times(fee);
        } else {
        	totalFee = amount.times(fee);
        }

        return totalFee.toFixed(8) + " " + feeCurrency + " (0" + stripZeroes(fee.times(100).toPrecision(2)) + "%)";
    }

    function totalMinusFee(order) {
        var fee = BigNumber(order.fee),
            total = BigNumber(order.total);
        if (order.type=="sell"){
        	total = total.minus(total.times(fee)).toFixed(8);
        } else {
        	total = total.toFixed(8);
        }
        return total;
    }

    function stripZeroes(text) {
        return /^0*(.*?)0*$/.exec(text)[1];
    }

    function makeTableRow(coin, market, obj) {
        var row = $(d.createElement('tr')).attr('data-id', obj.tradeID),
            coinNode = $(d.createElement('td')),
            typeNode = $(d.createElement('td')),
            categoryNode = $(d.createElement('td')),
            feeCurrency,
            formattedRate,
            formattedFee,
            formattedTotal;

        row.addClass('analysisRow selected');        
        typeNode.addClass(obj.type);

        coinNode.append(coin);
        if (market != "-")
        	coinNode.append($(d.createElement('span'))
                        	.text('/' + market));
        
        if (obj.type == "buy") typeNode.text('Buy');
        else if (obj.type == "sell") typeNode.text('Sell');

        if (obj.type == "buy") feeCurrency=coin;
        else if (obj.type == "sell") feeCurrency=market;
                
        var isLending = false;
        
        if (obj.category == "exchange") categoryNode.text('Exchange');
        else if (obj.category == "marginTrade") categoryNode.text('Margin Trade');
        else if (obj.category == "settlement") categoryNode.text('Settlement');
        else if (obj.category == "lendingSettlement"){
	        isLending = true;
	        categoryNode.text('Lending Fees');
	        obj.rate = "0";
	        formattedRate = "-";
	        formattedFee = "-";
	        formattedTotal = BigNumber(obj.total).toFixed(8);
	        typeNode.removeClass(obj.type);
	        typeNode.text('-');
	    } else if (obj.category == "lendingEarning"){
		    isLending = true;
		    typeNode.text('-');
		    categoryNode.html('<span class="buyClass">Lending Earning</span>');
		    formattedRate = (parseFloat(obj.rate)*100).toFixed(4) + "%";
		    formattedFee = (parseFloat(obj.fee)*100).toFixed(2) + "%";
		    formattedTotal = BigNumber(obj.total).toFixed(8);
	    }
	    
	    if (obj.category != "lendingSettlement" && obj.category != "lendingEarning"){
		    formattedFee = formatFee(obj,feeCurrency);
		    formattedTotal = totalMinusFee(obj);
		    formattedRate = BigNumber(obj.rate).toFixed(8);
		}

        row.append(makeTextCell("").addClass('indicator'));
        row.append(coinNode);
        row.append(typeNode);
        row.append(categoryNode);
        row.append(makeNumberCell(formattedRate));
        row.append(makeNumberCell(BigNumber(obj.amount).toFixed(8)));
        row.append(makeNumberCell(formattedFee));
        row.append(makeNumberCell(formattedTotal)
                   .append($(d.createElement('span'))
                           .text(" " + (isLending ? feeCurrency : market))));
        row.append(makeNumberCell(obj.date));

        return row;
    }

    function applyRowHandlers(row) {
        row.click(function() {
            var node = $(this),
                isSelected = node.hasClass('selected'),
                order = orderCache[node.attr('data-id')];

            node.toggleClass('selected', !isSelected);
            
            if (!isSelected) { // select
                addOrder(order);
                valueCache.selected++;
            } else { // deselect
                removeOrder(order);
                valueCache.selected--;
            }

            if (valueCache.selected != valueCache.totalOrders) {
                if (valueCache.selected == 0) {
                    $(".allIndicator").removeClass("mixed").removeClass("all");
                } else {
                    $(".allIndicator").removeClass("all").addClass("mixed");
                }
            } else {
                $(".allIndicator").removeClass("mixed").addClass("all");
            }

            refreshAnalysis();
        });
    }

    function noTradesAlert(orderBook) {
        var row = $(d.createElement('tr')),
            alert = $(d.createElement('td'))
                .attr('colspan', 8)
                .css('padding', '.8em'),
            p = $(d.createElement('p')).css('text-align','left');

        orderBook.append(row);
        row.append(alert);
        alert.append("<p><strong>No Trades</strong></p>");
        alert.append(p);
        p.append("You haven't made any trades yet. ");
        p.append("Once you have made some trades they will appear here.");
    }
        
    function noResultsAlert(orderBook) {
        var row = $(d.createElement('tr')),
            alert = $(d.createElement('td'))
                .attr('colspan', 8)
                .css('padding', '.8em'),
            link = $(d.createElement('a'))
                .attr('href', '#')
                .css('text-decoration', 'underline')
                .text("Clear All Filters"),
            p = $(d.createElement('p')).css('text-align','left');

        orderBook.append(row);
        row.append(alert);
        alert.append("<p><strong>No Results</strong></p>");
        alert.append(p);
        p.append("We didn't find any transactions within the specified search parameters. ");
        p.append("Try again or ");
        p.append(link);
        p.append(" to start over.");
        link.click(clearAnalysis);
    }

    function startOfDate(dt) {
        return moment.utc(dt).startOf('day').toDate();
    }

    function endOfDate(dt) {
        return moment.utc(dt).endOf('day').toDate();
    }

    function updateHash(stopEvent) {
        var coin = $("#coinName").val().trim().toUpperCase(),
            market = $("#marketName").val().trim().toUpperCase(),
            hash,
            range;

        if (coin == "" || market == "") {
        		range = getDateRange();
        		start = range.start;
        		end = range.end;
        		type = $("#showList").val().trim();
        		if (type=="sell")type=0;
        		if (type=="buy")type=1;
        		if (type=="all")type=2;
        		if (type=="lending")type=3;
        		populateTradeHistory(start,end,currentPage,tradesPerPage,type);
        		fetchTradeHistoryPages(start,end,tradesPerPage,type);
            return;
        }
        
        if ($("#showList").val().trim() == "lending"){
        	market = "-";
        } else if (market == "-"){
	        market = "BTC";
        }

        hash = "#analysis/" + 
               coin + '/' +
               market + '/' +
               $("#showList").val().trim()
       
        if (dateFilter != "All Dates") {
            range = getDateRange();
            hash += '/' + range.start + '/' + range.end;
        }

        if (stopEvent) {
            $(window)
                .off('hashchange', refreshHash)
                .one('hashchange', function(e) {
                    refreshHash(e, true);
                    $(window).on('hashchange', refreshHash);
                });
        }
        
        location.hash = hash;
    }

    function generateAnalysis() {
        var coin = $("#coinName").val().trim().toUpperCase(),
            market = $("#marketName").val().trim().toUpperCase(),
	        token = market == "-" ? coin : market + "_" + coin,
            typeFilter = $("#showList").val(),
            orderBook;
            
        // Remove original from DOM and store
        initialOrderBook.remove();

        // Remove old analysis if it's there.
        $(".orderBook tbody").remove();

        // Create empty table for us!
        orderBook = $(d.createElement('tbody')).appendTo($(".orderBook"));
   
        // Enable clear button
        $("#clearBtn").removeAttr('disabled');

        valueCache = {
            buyTotal: zero, // market currency total
            buyLendingTotal: zero,
            buyAmountTotal: zero, // amount sum
            sellTotal: zero,
            sellLendingTotal: zero,
            sellAmountTotal: zero,
            breakEvenTotal: zero,
            breakEvenAmountTotal: zero,
            buyFeeTotal: zero,
            buyFeeValueTotal: zero,
            buyFeeCount: 0,
            sellFeeTotal: zero,
            sellFeeValueTotal: zero,
            sellFeeCount: 0,
            market: market,
            coin: coin
        }
        
        $(".pages").hide();
        
        if (tradeHistory[token] == null) {
            refreshAnalysis();
            
            
            noResultsAlert(orderBook);
            
            
            return;
        }

        $(".calc").show();
        $(".info").hide();

        $(".orderBook").addClass("analysis");
        $(".tradeAnalysis").addClass("active");
        $(".allIndicator").removeClass("mixed").removeClass("all")

        tradeHistory[token].forEach(function(order) {
            var row;

			// Skip rows filtered out
			if (order.category == "lendingEarning" && typeFilter != 'lending'){
				return;
			} else if (order.category != "lendingEarning" && typeFilter != "all" && typeFilter != order.type){
                return;
			}
			
            row = makeTableRow(coin, market, order);
            orderBook.append(row);
            applyRowHandlers(row);
            
            addOrder(order);
            orderCache[order.tradeID] = order;
        });

        if (orderBook.children().length == 0) {
            noResultsAlert(orderBook);
        } else {
            valueCache.totalOrders = orderBook.children().length;
            valueCache.selected = valueCache.totalOrders;
            
            $(".allIndicator").addClass("all");
        }

        refreshAnalysis();
    };

    $("#analyzeBtn").click(updateHash);
    $("#showList").change(updateHash);

    $("#applyDateSelection").click(function() {
        //dateFilter = "User Defined";
        //applyDateRange();
        updateHash(true);
        $("#dateRangeSelector").fadeOut();
        lastDateFilter = dateFilter;
    });

    function disableControls(all) {
        $("#showList, #coinName, #marketName, " +
          "#dateRangeBtn, #analyzeBtn, #clearBtn").attr('disabled', true);
    }

    function enableControls() {
        $("#showList, #coinName, #marketName, " +
          "#dateRangeBtn, #analyzeBtn").removeAttr('disabled');

        if ($(".orderBook").hasClass('analysis')) {
            $("#clearBtn").removeAttr('disabled');
        }
    }

    function setWidgetDate(node, v) {
        // Workaround select event triggering on programmatical set

		// Adjust for timezone
        var offset0 = v[0].getTimezoneOffset() * 60 * 1000;
        var offset1 = v[1].getTimezoneOffset() * 60 * 1000;
		var t0 = v[0].getTime();
		var t1 = v[1].getTime();
		v[0] = new Date(t0 + offset0);
		v[1] = new Date(t1 + offset1);
		
        var oldSelect = node.datepick('option', 'onSelect');
        node.datepick('option', 'onSelect', function(){});
        node.datepick('setDate', v);
        node.datepick('option', 'onSelect', oldSelect);
    }

    function refreshHash(e, noDateCheck) {
        var hashData = parseHash(location.hash),
            range;

        if (hashData != null && hashData.coin && hashData.market) {
            $("#coinName").val(hashData.coin);
            $("#marketName").val(hashData.market);

            refreshAnalyzeControls();
            
            if (hashData.market == "-"){
	            $("#showList").val("lending");
            } else if (hashData.show) {
                $("#showList").val(hashData.show);
            }
            
            if (!noDateCheck && hashData.fromDate && hashData.toDate) {
                dateFilter = "User Defined";
                
                setWidgetDate($("#dateWidget"), [
                    new Date(hashData.fromDate * 1000),
                    new Date(hashData.toDate * 1000)
                ]);

                applyDateRange();
            }

            range = getDateRange();

            if (lastStart == range.start && lastEnd == range.end) {
                generateAnalysis();
                return;
            }

            lastStart = range.start;
            lastEnd = range.end;

            // Get new data if the date changed.
            disableControls();
            renderTradeHistoryPages(0);
            $(".orderBook").hide()
                .after($("<div style='text-align: center; margin-top: 1em; font-size: 200%'>" +
                         "Loading...</div>"));
            $.get('/private.php', {command: 'returnPersonalTradeHistory',
                                   start: range.start,
                                   end: range.end}).done(function(data) {
                try {
                    tradeHistory = JSON.parse(data);
                } catch (e) {
                    // Not authed probably due to session death. Back to login!
                    location.href = "/login";
                    return;
                }
                generateAnalysis();

                $(".orderBook").show().next().remove();
                enableControls();
            });
        }
    }

    $(window).on('hashchange', refreshHash);
    refreshAnalyzeControls();
    refreshHash();
})();
