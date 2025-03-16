from pywps import WPS
wps = WPS()
wps.open("example.et")
wps.save_as("example.xlsx", file_format=WPS.FileFormat.XLSX)