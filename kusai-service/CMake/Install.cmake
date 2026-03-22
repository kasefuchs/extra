include(GNUInstallDirs)

install(
  TARGETS kusai-service
  RUNTIME DESTINATION ${CMAKE_INSTALL_BINDIR}
)

include(CPack)
