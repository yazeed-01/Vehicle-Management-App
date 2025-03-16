import 'package:car/api/api_service.dart';
import 'package:flutter/material.dart';

class MyVehiclesPage extends StatefulWidget {
  @override
  _MyVehiclesPageState createState() => _MyVehiclesPageState();
}

class _MyVehiclesPageState extends State<MyVehiclesPage> {
  final _apiService = ApiService();
  late Future<List<Map<String, dynamic>>> _vehiclesFuture;

  @override
  void initState() {
    super.initState();
    _vehiclesFuture = _apiService.getUserVehicles();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('My Vehicles'),
      ),
      body: FutureBuilder<List<Map<String, dynamic>>>(
        future: _vehiclesFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return Center(child: CircularProgressIndicator());
          } else if (snapshot.hasError) {
            return Center(child: Text('Error: ${snapshot.error}'));
          } else if (!snapshot.hasData || snapshot.data!.isEmpty) {
            return Center(child: Text('No vehicles found'));
          }

          final vehicles = snapshot.data!;
          return ListView.builder(
            padding: EdgeInsets.all(16),
            itemCount: vehicles.length,
            itemBuilder: (context, index) {
              final vehicle = vehicles[index];
              return Card(
                elevation: 4,
                margin: EdgeInsets.only(bottom: 16),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Padding(
                  padding: EdgeInsets.all(16),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        '${vehicle['model']}',
                        style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      SizedBox(height: 8),
                      Text('Plate: ${vehicle['car_plate']}'),
                      Text('Color: ${vehicle['color']}'),
                      Text('Type: ${vehicle['type']}'),
                    ],
                  ),
                ),
              );
            },
          );
        },
      ),
    );
  }
}
