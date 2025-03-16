import 'package:car/api/api_service.dart';
import 'package:flutter/material.dart';
import 'vehicle_details_page.dart';

class SearchPlatePage extends StatefulWidget {
  @override
  _SearchPlatePageState createState() => _SearchPlatePageState();
}

class _SearchPlatePageState extends State<SearchPlatePage> {
  final _formKey = GlobalKey<FormState>();
  String _plateNumber = '';
  final _apiService = ApiService();
  bool _isLoading = false;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Search by Plate'),
      ),
      body: Padding(
        padding: EdgeInsets.all(16),
        child: Form(
          key: _formKey,
          child: Column(
            children: [
              Card(
                elevation: 4,
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Padding(
                  padding: EdgeInsets.all(16),
                  child: TextFormField(
                    decoration: InputDecoration(
                      labelText: 'Plate Number',
                      prefixIcon: Icon(Icons.directions_car),
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                    ),
                    validator: (value) {
                      if (value == null || value.isEmpty) {
                        return 'Please enter a plate number';
                      }
                      return null;
                    },
                    onSaved: (value) {
                      _plateNumber = value!;
                    },
                  ),
                ),
              ),
              SizedBox(height: 16),
              _isLoading
                  ? CircularProgressIndicator()
                  : ElevatedButton(
                      style: ElevatedButton.styleFrom(
                        minimumSize: Size(double.infinity, 50),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(12),
                        ),
                      ),
                      onPressed: () async {
                        if (_formKey.currentState!.validate()) {
                          _formKey.currentState!.save();
                          setState(() {
                            _isLoading = true;
                          });
                          try {
                            final result =
                                await _apiService.searchByPlate(_plateNumber);
                            Navigator.push(
                              context,
                              MaterialPageRoute(
                                builder: (context) => VehicleDetailsPage(
                                  vehicleData: result,
                                ),
                              ),
                            );
                          } catch (e) {
                            ScaffoldMessenger.of(context).showSnackBar(
                              SnackBar(
                                content: Text('Search failed: $e'),
                                backgroundColor: Colors.red,
                              ),
                            );
                          } finally {
                            setState(() {
                              _isLoading = false;
                            });
                          }
                        }
                      },
                      child: Text(
                        'Search',
                        style: TextStyle(fontSize: 18),
                      ),
                    ),
            ],
          ),
        ),
      ),
    );
  }
}
