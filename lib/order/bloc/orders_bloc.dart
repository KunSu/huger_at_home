import 'dart:async';

import 'package:authentication_repository/authentication_repository.dart';
import 'package:bloc/bloc.dart';
import 'package:equatable/equatable.dart';
import 'package:fe/address/addresses_repository.dart';
import 'package:fe/order/order.dart';
import 'package:meta/meta.dart';

part 'orders_event.dart';
part 'orders_state.dart';

class OrdersBloc extends Bloc<OrdersEvent, OrdersState> {
  OrdersBloc({
    @required this.authenticationRepository,
    @required this.addressesRepository,
    @required this.ordersRepository,
  }) : super(OrdersLoadInProgress());

  final AuthenticationRepository authenticationRepository;
  final AddressesRepository addressesRepository;
  final OrdersRepository ordersRepository;

  @override
  Stream<OrdersState> mapEventToState(OrdersEvent event) async* {
    if (event is OrdersLoaded) {
      yield* _mapOrdersLoadedToState();
    } else if (event is OrderAdded) {
      yield* _mapOrderAddedToState(event);
    } else if (event is OrderUpdated) {
      yield* _mapOrderUpdatedToState(event);
    } else if (event is OrderDeleted) {
      yield* _mapOrderDeletedToState(event);
    }
  }

  Stream<OrdersState> _mapOrdersLoadedToState() async* {
    try {
      await ordersRepository.reload(userID: authenticationRepository.user.id);
      final orders = ordersRepository.loadOrders();
      yield OrdersLoadSuccess(
        orders,
      );
    } catch (e) {
      yield OrdersLoadFailure(e.toString());
    }
  }

  Stream<OrdersState> _mapOrderAddedToState(OrderAdded event) async* {
    if (state is OrdersLoadSuccess) {
      try {
        final List<Order> updatedOrders =
            List.from((state as OrdersLoadSuccess).orders)..add(event.order);
        yield OrdersLoadSuccess(updatedOrders);
        _saveOrders(updatedOrders);
      } catch (_) {
        yield const OrdersLoadFailure('Internal Error');
      }
    }
  }

  Stream<OrdersState> _mapOrderUpdatedToState(OrderUpdated event) async* {
    if (state is OrdersLoadSuccess) {
      final List<Order> updatedOrders =
          (state as OrdersLoadSuccess).orders.map((order) {
        return order.id == event.order.id ? event.order : order;
      }).toList();
      yield OrdersLoadSuccess(updatedOrders);
      _saveOrders(updatedOrders);
    }
  }

  Stream<OrdersState> _mapOrderDeletedToState(OrderDeleted event) async* {
    if (state is OrdersLoadSuccess) {
      final updatedOrders = (state as OrdersLoadSuccess)
          .orders
          .where((order) => order.id != event.order.id)
          .toList();
      yield OrdersLoadSuccess(updatedOrders);
      _saveOrders(updatedOrders);
    }
  }

  Future _saveOrders(List<Order> orders) {
    return ordersRepository.saveOrders(
      orders,
    );
  }
}
