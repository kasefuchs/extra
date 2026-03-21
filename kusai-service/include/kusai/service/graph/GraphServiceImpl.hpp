#pragma once

#include "kusai/service/proto/graph/v1.grpc.pb.h"

namespace kusai::service {
class GraphServiceImpl final : proto::graph::v1::GraphService::Service {};
}  // namespace kusai::service
