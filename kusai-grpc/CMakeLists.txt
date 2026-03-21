cmake_minimum_required(VERSION 3.16)

project(kusai-grpc
  VERSION 0.0.0
  LANGUAGES CXX
)

list(APPEND CMAKE_MODULE_PATH "${CMAKE_CURRENT_SOURCE_DIR}/CMake")

include(Settings)
include(Dependencies)

add_subdirectory(proto)
add_subdirectory(src)
add_subdirectory(app)
