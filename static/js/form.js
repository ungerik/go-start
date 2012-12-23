var gostart_form = gostart_form || {};

gostart_form.swapChildren = function (a, b) {
	var x = a.children().detach();
	var y = b.children().detach();
	x.appendTo(b);
	y.appendTo(a);
};

gostart_form.swapAttribs = function (a, b, attr) {
	var x = a.attr(attr);
	var y = b.attr(attr);
	a.attr(attr, y);
	b.attr(attr, x);
};

gostart_form.swapValues = function(a, b) {
	var x = a.val();
	var y = b.val();
	a.val(y);
	b.val(x);
};

gostart_form.swapChecked = function(a, b) {
	var x = a.prop("checked");
	var y = b.prop("checked");
	a.prop("checked", y);
	b.prop("checked", x);
};

gostart_form.swapRowValues = function(tr0, tr1) {
	var inputs0 = tr0.find("td :input").not(":button");
	var inputs1 = tr1.find("td :input").not(":button");
	for (i=0; i < inputs0.length; i++) {
		gostart_form.swapValues(inputs0.eq(i), inputs1.eq(i));
	}
	inputs0 = tr0.find("td :checkbox");
	inputs1 = tr1.find("td :checkbox");
	for (i=0; i < inputs0.length; i++) {
		gostart_form.swapChecked(inputs0.eq(i), inputs1.eq(i));
	}
	// Todo: Generic solution for special case of media.ImageRef
	var img0 = tr0.find("td .thumbnail-frame");
	var img1 = tr1.find("td .thumbnail-frame");
	for (i=0; i < img0.length; i++) {
		gostart_form.swapChildren(img0.eq(i), img1.eq(i));
	}
};

gostart_form.onLengthChanged = function(table) {
	var rows = table.find("tr");
	table.prev("input[type=hidden]").val(rows.length-1);
	rows.each(function(row) {
		var firstRow = (row == 1); // ignore header row
		var lastRow = (row == rows.length-1);
		var buttons = jQuery(this).find("td:last > :button");
		buttons.eq(0).prop("disabled", firstRow);
		buttons.eq(1).prop("disabled", lastRow);
		if (lastRow) {
			buttons.eq(2).attr("onclick", "gostart_form.addRow(this);").text("+");
		} else {
			buttons.eq(2).attr("onclick", "gostart_form.removeRow(this);").html("&times;");
		}
	});
};

gostart_form.removeRow = function(button) {
	if (confirm("Are you sure you want to delete this row?")) {
		var tr = jQuery(button).parents("tr");
		var table = tr.parents("table");

		// Swap all values with following rows to move the values of the
		// row to be deleted to the last row and everything else one row up
		var rows = tr.add(tr.nextAll());
		rows.each(function(i) {
			if (i > 0)
				gostart_form.swapRowValues(rows.eq(i-1), rows.eq(i));
		});

		rows.last().remove();

		gostart_form.onLengthChanged(table);
	}
};

gostart_form.addRow = function(button) {
	var tr0 = jQuery(button).parents("tr");
	var tr1 = tr0.clone();
	var table = tr0.parents("table");

	// Set correct class for new row
	var numRows = tr0.prevAll().length + 1;
	var evenOdd = (numRows % 2 === 0) ? " even" : " odd";
	tr1.attr("class", "row"+numRows+evenOdd);

	// Correct name attributes of the new row's input elements
	var oldIndex = "."+(numRows-2);
	var newIndex = "."+(numRows-1);
	tr1.find("td :input").not(":button").each(function() {
		var i = this.name.lastIndexOf(oldIndex);
		this.name = this.name.slice(0,i)+newIndex+this.name.slice(i+oldIndex.length);
	});

	// Make IDs of new row content unique by adding "_"
	var ids = [];
	tr1.find("[id]").each(function() {
		var elem = jQuery(this);
		var id = elem.attr("id");
		ids.push(id);
		elem.attr("id", id+"_");
	});
	tr1.find("[onclick]").each(function() {
		var elem = jQuery(this);
		var onclick = elem.attr("onclick");
		for (var i = 0; i < ids.length; i++) {
			onclick = onclick.replace("'#"+ids[i]+"'", "'#"+ids[i]+"_'");
			onclick = onclick.replace('"#'+ids[i]+'"', '"#'+ids[i]+'_"');
		}
		elem.attr("onclick", onclick);
	});

	tr1.insertAfter(tr0);

	gostart_form.onLengthChanged(table);
};

gostart_form.moveRowUp = function(button) {
	var tr1 = jQuery(button).parents("tr");
	var tr0 = tr1.prev();
	gostart_form.swapRowValues(tr0, tr1);
};

gostart_form.moveRowDown = function(button) {
	var tr0 = jQuery(button).parents("tr");
	var tr1 = tr0.next();
	gostart_form.swapRowValues(tr0, tr1);
};

