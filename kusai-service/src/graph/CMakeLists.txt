add_library(graph
  GraphServiceImpl.cpp
)

target_include_directories(graph PUBLIC
  ${CMAKE_SOURCE_DIR}/include
)

target_sources(graph PUBLIC
  FILE_SET HEADERS
  BASE_DIRS ${CMAKE_SOURCE_DIR}/include
  FILES
  ${CMAKE_SOURCE_DIR}/include/kusai/service/graph/GraphServiceImpl.hpp
)

target_link_libraries(graph
  PUBLIC
    proto
)
