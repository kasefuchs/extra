add_library(textchain
  TextChainServiceImpl.cpp
)

target_include_directories(textchain PUBLIC
  ${CMAKE_SOURCE_DIR}/include
)

target_sources(textchain PUBLIC
  FILE_SET HEADERS
  BASE_DIRS ${CMAKE_SOURCE_DIR}/include
  FILES
  ${CMAKE_SOURCE_DIR}/include/kusai/service/textchain/TextChainServiceImpl.hpp
)

target_link_libraries(textchain
  PUBLIC
    proto
    registry
    kusai::textchain
)
