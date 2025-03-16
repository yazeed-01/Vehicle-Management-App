import 'package:car/api/api_service.dart';
import 'package:car/auth/login.dart';
import 'package:flutter/material.dart';
import 'home_page.dart';

void main() {
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  final ApiService _apiService = ApiService();

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        brightness: Brightness.light,
        primaryColor: Colors.blue,
      ),
      darkTheme: ThemeData(
        brightness: Brightness.dark,
        primaryColor: Colors.blue[900],
        scaffoldBackgroundColor: Colors.grey[900],
        cardColor: Colors.grey[850],
      ),
      themeMode: ThemeMode.system,
      home: FutureBuilder<bool>(
        future: _apiService.isLoggedIn(),
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return Scaffold(
              body: Center(child: CircularProgressIndicator()),
            );
          }
          if (snapshot.hasData && snapshot.data == true) {
            return HomePage();
          }
          return LoginPage();
        },
      ),
    );
  }
}
