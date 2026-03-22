#include "kusai/service/graph/GraphServiceImpl.hpp"

#include <kusai/graph/MemoryGraph.hpp>

namespace kusai::service {
grpc::Status GraphServiceImpl::Create(grpc::ServerContext* context, const proto::graph::v1::CreateRequest* request,
                                      proto::graph::v1::CreateResponse* response) {
  std::shared_ptr<AbstractGraph> graph;
  switch (request->impl_case()) {
    case proto::graph::v1::CreateRequest::kMemory:
      graph = std::make_shared<MemoryGraph>();
      break;
    default:
      return {grpc::StatusCode::INVALID_ARGUMENT, "Unknown implementation"};
  }

  const auto id = graphRegistry_->add(graph);
  response->set_graph_id(id);

  return grpc::Status::OK;
}

grpc::Status GraphServiceImpl::Delete(grpc::ServerContext* context, const proto::graph::v1::DeleteRequest* request,
                                      proto::graph::v1::DeleteResponse* response) {
  const auto success = graphRegistry_->remove(request->graph_id());
  response->set_success(success);

  return grpc::Status::OK;
}
}  // namespace kusai::service
