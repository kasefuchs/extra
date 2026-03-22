#pragma once

#include <kusai/graph/AbstractGraph.hpp>

#include "kusai/service/proto/graph/v1.grpc.pb.h"
#include "kusai/service/registry/Registry.hpp"

namespace kusai::service {
class GraphServiceImpl final : public proto::graph::v1::GraphService::Service {
 public:
  explicit GraphServiceImpl(const std::shared_ptr<Registry<AbstractGraph>>& graphRegistry)
      : graphRegistry_(graphRegistry) {}

  ~GraphServiceImpl() override = default;

  grpc::Status Create(grpc::ServerContext* context, const proto::graph::v1::CreateRequest* request,
                      proto::graph::v1::CreateResponse* response) override;

  grpc::Status Delete(grpc::ServerContext* context, const proto::graph::v1::DeleteRequest* request,
                      proto::graph::v1::DeleteResponse* response) override;

 private:
  std::shared_ptr<Registry<AbstractGraph>> graphRegistry_;
};
}  // namespace kusai::service
