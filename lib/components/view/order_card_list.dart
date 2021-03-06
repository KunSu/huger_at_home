import 'package:fe/components/ult/status_color.dart';
import 'package:fe/components/view/order/order_action.dart';
import 'package:fe/order/order.dart';
import 'package:flutter/material.dart';
import 'package:flutter_form_bloc/flutter_form_bloc.dart';

class OrderCardList extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return BlocBuilder<OrdersBloc, OrdersState>(
      builder: (context, state) {
        if (state is OrdersLoadInProgress) {
          return const CircularProgressIndicator();
        }
        if (state is OrdersLoadSuccess) {
          if (state.orders == null || state.orders.isEmpty) {
            return const Text('You do not have any order yet');
          }
          return ListView.builder(
            itemCount: state.orders.length,
            itemBuilder: (context, index) =>
                OrderCardView(order: state.orders[index]),
          );
        }
        return const Text('Something went wrong!');
      },
    );
  }
}

class OrderCardView extends StatelessWidget {
  const OrderCardView({Key key, this.order}) : super(key: key);
  final Order order;
  @override
  Widget build(BuildContext context) {
    return Card(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          ListTile(
            title: Text('Order #: ${order.id}'),
            subtitle: RichText(
              text: TextSpan(
                text: 'Pick up date: ${order.pickupDateAndTime}\n',
                style: DefaultTextStyle.of(context).style,
                children: <TextSpan>[
                  TextSpan(
                    text: 'Address: ${order.address}\n',
                  ),
                  const TextSpan(
                    text: 'Status: ',
                  ),
                  TextSpan(
                    text: '${order.status}\n',
                    style: TextStyle(
                      fontWeight: FontWeight.bold,
                      color: getStatusColor(status: order.status),
                    ),
                  ),
                ],
              ),
            ),
          ),
          OrderAction(
            order: order,
          ),
        ],
      ),
    );
  }
}
