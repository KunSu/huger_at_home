import 'package:authentication_repository/authentication_repository.dart';
import 'package:fe/components/view/buttom_navigation_bar.dart';
import 'package:fe/order/models/model.dart';
import 'package:fe/pantry/models/appbar.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:fe/cart/cart.dart';
import 'package:fe/pantry/pantry.dart';

class PantryPage extends StatefulWidget {
  static String routeName = '/pantry';

  @override
  _PantryPageState createState() => _PantryPageState();
}

class _PantryPageState extends State<PantryPage> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: CustomScrollView(
        slivers: [
          const MyAppBar(
            title: 'Pantry',
          ),
          const SliverToBoxAdapter(child: SizedBox(height: 12)),
          BlocBuilder<CatalogBloc, CatalogState>(
            builder: (context, state) {
              if (state is CatalogLoading) {
                return const SliverFillRemaining(
                  child: Center(
                    child: CircularProgressIndicator(),
                  ),
                );
              }
              if (state is CatalogLoaded) {
                return SliverList(
                  delegate: SliverChildBuilderDelegate(
                    (context, index) => _MyListItem(
                      state.catalog.getByPosition(index),
                    ),
                  ),
                );
              }

              return const Text('Something went wrong!');
            },
          ),
        ],
      ),
      bottomNavigationBar: MyBottomNavigationBar(
        identity: RepositoryProvider.of<AuthenticationRepository>(context)
            .user
            .userIdentity,
      ),
    );
  }
}

class _AddButton extends StatelessWidget {
  const _AddButton({Key key, @required this.item}) : super(key: key);

  final Item item;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return BlocBuilder<CartBloc, CartState>(
      builder: (context, state) {
        if (state is CartLoading) {
          return const CircularProgressIndicator();
        }
        if (state is CartLoaded) {
          return FlatButton(
            onPressed: state.cart.items.contains(item)
                ? null
                : () => context.read<CartBloc>().add(CartItemAdded(item)),
            splashColor: theme.primaryColor,
            child: state.cart.items.contains(item)
                ? const Icon(Icons.check, semanticLabel: 'ADDED')
                : const Text('ADD'),
          );
        }
        return const Text('Something went wrong!');
      },
    );
  }
}

class _MyListItem extends StatelessWidget {
  const _MyListItem(this.item, {Key key}) : super(key: key);

  final Item item;

  @override
  Widget build(BuildContext context) {
    final textTheme = Theme.of(context).textTheme.headline6;
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      child: LimitedBox(
        maxHeight: 48,
        child: Row(
          children: [
            const SizedBox(width: 24),
            Expanded(child: Text(item.name, style: textTheme)),
            const SizedBox(width: 24),
            // Text(item.pickupDateAndTime),
            Text('Quantity' + item.quantityNumber),
            Text(item.quantityUnit),
            // _AddButton(item: item),
          ],
        ),
      ),
    );
  }
}
