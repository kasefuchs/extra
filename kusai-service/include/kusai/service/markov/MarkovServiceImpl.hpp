#pragma once

#include "kusai/service/proto/markov/v1.grpc.pb.h"

namespace kusai::service {
class MarkovServiceImpl final : proto::markov::v1::MarkovService::Service {};
}  // namespace kusai::service
