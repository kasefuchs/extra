add_library(registry INTERFACE)

target_include_directories(registry INTERFACE
  ${CMAKE_SOURCE_DIR}/include
)

target_sources(registry INTERFACE
  FILE_SET HEADERS
    BASE_DIRS ${CMAKE_SOURCE_DIR}/include
    FILES
      ${CMAKE_SOURCE_DIR}/include/kusai/service/registry/Registry.hpp
)
