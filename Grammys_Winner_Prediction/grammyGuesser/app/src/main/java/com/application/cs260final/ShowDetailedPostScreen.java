package com.application.cs260final;

import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.widget.ImageView;
import android.widget.Toast;

import androidx.appcompat.app.AppCompatActivity;

public class ShowDetailedPostScreen extends AppCompatActivity {
    ImageView profile_button1;
    ImageView search_button;
    ImageView predictions_button;
    ImageView notification_button;
    ImageView profile_button2;
    ImageView back_button;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_home_screen_click_on_post);

        //1. CREATE AN ACTION LINKED TO THE PROFILE BUTTON BUTTON (AT TOP)
        //-> setOnClickListener is listening if button is ever clicked. If it is, all the actions in this function run.
        //TODO: ALL MENU BUTTONS ARE GROUPED AS ONE. CANNOT GET SEARCH BUTTON'S ID
        profile_button1 = (ImageView) findViewById(R.id.imageProfile1);
        profile_button1.setOnClickListener(new View.OnClickListener(){

            //When button is clicked, run code inside
            public void onClick(View v) {
                //TODO: CAN REMOVE THIS AFTER EVERYTHING WORKS CORRECTLY
                //specifies what class its in, what the text says, how long the text will be there, and if we should show it or not.
                Toast.makeText(ShowDetailedPostScreen.this, "Going to Profile Page", Toast.LENGTH_SHORT).show();

                //Links the action of pressing button = going to next layout.
                Intent intent = new Intent(ShowDetailedPostScreen.this, Profile.class);

                //Start Activity on a new screen that's encoded in the activity_home_screen_click_on_post class file.
                startActivity(intent);
            }
        });

        //implement the back button even tho we prob don't need it
        back_button = findViewById(R.id.imageBackButton);
        back_button.setOnClickListener(new View.OnClickListener() {
            public void onClick(View v) {
                Intent intent = new Intent(ShowDetailedPostScreen.this, HomeScreenActivity.class);
                startActivity(intent);
            }
        });

  // ************************************* MENU BUTTON FUNCTIONALITIES BELOW **************************************** //

        //2. CREATE AN ACTION LINKED TO THE SEARCH BUTTON (AT THE MENU BAR BELOW)
        //-> setOnClickListener is listening if button is ever clicked. If it is, all the actions in this function run.
        //TODO: ALL MENU BUTTONS ARE GROUPED AS ONE. CANNOT GET SEARCH BUTTON'S ID
        search_button = (ImageView) findViewById(R.id.imageSearch);
        search_button.setOnClickListener(new View.OnClickListener(){

            //When button is clicked, run code inside
            public void onClick(View v) {
                //TODO: CAN REMOVE THIS AFTER EVERYTHING WORKS CORRECTLY
                //specifies what class its in, what the text says, how long the text will be there, and if we should show it or not.
                Toast.makeText(ShowDetailedPostScreen.this, "Going to Search Page", Toast.LENGTH_SHORT).show();

                //Links the action of pressing button = going to next layout.
                Intent intent = new Intent(ShowDetailedPostScreen.this, Search.class);

                //Start Activity on a new screen that's encoded in the activity_home_screen_click_on_post class file.
                startActivity(intent);
            }
        });

        //3. CREATE AN ACTION LINKED TO THE PREDITCION BUTTON (AT THE MENU BAR BELOW)
        //-> setOnClickListener is listening if button is ever clicked. If it is, all the actions in this function run.
        //TODO: ALL MENU BUTTONS ARE GROUPED AS ONE. CANNOT GET SEARCH BUTTON'S ID
        predictions_button = (ImageView) findViewById(R.id.imageVote);
        predictions_button.setOnClickListener(new View.OnClickListener(){

            //When button is clicked, run code inside
            public void onClick(View v) {
                //TODO: CAN REMOVE THIS AFTER EVERYTHING WORKS CORRECTLY
                //specifies what class its in, what the text says, how long the text will be there, and if we should show it or not.
                Toast.makeText(ShowDetailedPostScreen.this, "Going to Predictions Page", Toast.LENGTH_SHORT).show();

                //Links the action of pressing button = going to next layout.
                Intent intent = new Intent(ShowDetailedPostScreen.this, Predictions.class);

                //Start Activity on a new screen that's encoded in the activity_home_screen_click_on_post class file.
                startActivity(intent);
            }
        });

        //4. CREATE AN ACTION LINKED TO THE NOTIFICATIONS BUTTON (AT THE MENU BAR BELOW)
        //-> setOnClickListener is listening if button is ever clicked. If it is, all the actions in this function run.
        //TODO: ALL MENU BUTTONS ARE GROUPED AS ONE. CANNOT GET SEARCH BUTTON'S ID
        notification_button = (ImageView) findViewById(R.id.imageBell);
        notification_button.setOnClickListener(new View.OnClickListener(){

            //When button is clicked, run code inside
            public void onClick(View v) {
                //TODO: CAN REMOVE THIS AFTER EVERYTHING WORKS CORRECTLY
                //specifies what class its in, what the text says, how long the text will be there, and if we should show it or not.
                Toast.makeText(ShowDetailedPostScreen.this, "Going to Notifications Page", Toast.LENGTH_SHORT).show();

                //Links the action of pressing button = going to next layout.
                Intent intent = new Intent(ShowDetailedPostScreen.this, Notifications.class);

                //Start Activity on a new screen that's encoded in the activity_home_screen_click_on_post class file.
                startActivity(intent);
            }
        });

        //5. CREATE AN ACTION LINKED TO THE PROFILE BUTTON (AT THE MENU BAR BELOW)
        //-> setOnClickListener is listening if button is ever clicked. If it is, all the actions in this function run.
        //TODO: ALL MENU BUTTONS ARE GROUPED AS ONE. CANNOT GET SEARCH BUTTON'S ID
        profile_button2 = (ImageView) findViewById(R.id.imageProfileButton);
        profile_button2.setOnClickListener(new View.OnClickListener(){

            //When button is clicked, run code inside
            public void onClick(View v) {
                //TODO: CAN REMOVE THIS AFTER EVERYTHING WORKS CORRECTLY
                //specifies what class its in, what the text says, how long the text will be there, and if we should show it or not.
                Toast.makeText(ShowDetailedPostScreen.this, "Going to Profile Page", Toast.LENGTH_SHORT).show();

                //Links the action of pressing button = going to next layout.
                Intent intent = new Intent(ShowDetailedPostScreen.this, Profile.class);

                //Start Activity on a new screen that's encoded in the activity_home_screen_click_on_post class file.
                startActivity(intent);
            }
        });

    }
}
