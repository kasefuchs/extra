add_library(proto
  kusai/service/proto/graph/v1.proto
  kusai/service/proto/markov/v1.proto
  kusai/service/proto/textchain/v1.proto
  kusai/service/proto/tokenizer/v1.proto
)

protobuf_generate(
  TARGET proto
  LANGUAGE cpp
)

protobuf_generate(
  TARGET proto
  LANGUAGE grpc
  PLUGIN "protoc-gen-grpc=\$<TARGET_FILE:gRPC::grpc_cpp_plugin>"
  GENERATE_EXTENSIONS .grpc.pb.h .grpc.pb.cc
)

target_include_directories(proto
  PUBLIC
    ${CMAKE_CURRENT_BINARY_DIR}
    ${CMAKE_CURRENT_SOURCE_DIR}
)

target_link_libraries(proto
  PUBLIC
    gRPC::grpc++
    protobuf::libprotobuf
)
