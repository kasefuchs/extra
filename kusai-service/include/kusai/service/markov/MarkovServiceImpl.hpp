#pragma once

#include <kusai/markov/AbstractMarkov.hpp>

#include "kusai/service/proto/markov/v1.grpc.pb.h"
#include "kusai/service/registry/Registry.hpp"

namespace kusai::service {
class MarkovServiceImpl final : public proto::markov::v1::MarkovService::Service {
 public:
  explicit MarkovServiceImpl(const std::shared_ptr<Registry<AbstractGraph>>& graphRegistry,
                             const std::shared_ptr<Registry<AbstractMarkov>>& markovRegistry)
      : graphRegistry_(graphRegistry), markovRegistry_(markovRegistry) {}

  ~MarkovServiceImpl() override = default;

  grpc::Status Create(grpc::ServerContext* context, const proto::markov::v1::CreateRequest* request,
                      proto::markov::v1::CreateResponse* response) override;

  grpc::Status Delete(grpc::ServerContext* context, const proto::markov::v1::DeleteRequest* request,
                      proto::markov::v1::DeleteResponse* response) override;

 private:
  std::shared_ptr<Registry<AbstractGraph>> graphRegistry_;
  std::shared_ptr<Registry<AbstractMarkov>> markovRegistry_;
};
}  // namespace kusai::service
