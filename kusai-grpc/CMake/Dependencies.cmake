# Protobuf
find_package(protobuf CONFIG QUIET)
if (NOT protobuf_FOUND)
  find_package(Protobuf REQUIRED)
endif()

# gRPC
find_package(gRPC CONFIG REQUIRED)

# KusAI
find_package(kusai CONFIG REQUIRED)
