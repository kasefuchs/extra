#include "kusai/service/markov/MarkovServiceImpl.hpp"

#include <kusai/markov/BackoffMarkov.hpp>
#include <kusai/markov/NGramMarkov.hpp>
#include <kusai/markov/SimpleMarkov.hpp>

namespace kusai::service {
grpc::Status MarkovServiceImpl::Create(grpc::ServerContext* context, const proto::markov::v1::CreateRequest* request,
                                       proto::markov::v1::CreateResponse* response) {
  const auto graph = graphRegistry_->get(request->graph_id());
  if (!graph) return {grpc::StatusCode::NOT_FOUND, "Graph not found"};

  std::shared_ptr<AbstractMarkov> markov;
  switch (request->impl_case()) {
    case proto::markov::v1::CreateRequest::kSimple:
      markov = std::make_shared<SimpleMarkov>(*graph);
      break;
    case proto::markov::v1::CreateRequest::kNgram:
      markov = std::make_shared<NGramMarkov>(*graph, request->ngram().context_size());
      break;
    case proto::markov::v1::CreateRequest::kBackoff:
      markov = std::make_shared<BackoffMarkov>(*graph, request->backoff().max_context_size());
      break;
    default:
      return {grpc::StatusCode::INVALID_ARGUMENT, "Unknown implementation"};
  }

  const auto id = markovRegistry_->add(markov);
  response->set_markov_id(id);

  return grpc::Status::OK;
}

grpc::Status MarkovServiceImpl::Delete(grpc::ServerContext* context, const proto::markov::v1::DeleteRequest* request,
                                       proto::markov::v1::DeleteResponse* response) {
  const auto success = markovRegistry_->remove(request->markov_id());
  response->set_success(success);

  return grpc::Status::OK;
}
}  // namespace kusai::service
