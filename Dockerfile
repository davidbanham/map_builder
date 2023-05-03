FROM osgeo/gdal

ADD map_builder map_builder

ENTRYPOINT ["/map_builder"]
