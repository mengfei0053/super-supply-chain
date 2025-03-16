import XlsxTemplate from "xlsx-template";
import path from "path";
import fs from "fs";

fs.readFile("./invoice_tmp.xlsx", async function (err, data) {
  // Create a template
  var template = new XlsxTemplate(data);

  // Replacements take place on first sheet
  var sheetNumber = 1;

  // Set up some placeholder values matching the placeholders in the template
  var values = {
    extractDate: new Date(),
    dates: [
      new Date("2013-06-01"),
      new Date("2013-06-02"),
      new Date("2013-06-03"),
    ],
    people: [
      { name: "John Smith", age: 20 },
      { name: "Bob Johnson", age: 22 },
    ],
  };

  // Perform substitution
  template.substitute(sheetNumber, values);
  template.substitute("2-发票明细信息", values);

  // Get binary data
  var data = template.generate();

  fs.writeFileSync("invoice.xlsx", data, "binary");
});
